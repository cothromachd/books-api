package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Book struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var	user  = flag.String("u", "username", "Username to connect database\npostgres://{username}:{password}@{hostname}:{port}/{database-name}")
var	password  = flag.String("pw", "password", "Password to connect database\npostgres://{username}:{password}@{hostname}:{port}/{database-name}")
var hostname  = flag.String("hn", "hostname", "Hostname to connect database\npostgres://{username}:{password}@{hostname}:{port}/{database-name}")
var	port  = flag.String("p", "5432", "Port to connect database\npostgres://{username}:{password}@{hostname}:{port}/{database-name}")
var	database  = flag.String("db", "database-name", "Database name to connect\npostgres://{username}:{password}@{hostname}:{port}/{database-name}")


func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		},
	})

	flag.Parse()
	dbPool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s", *user, *password, *hostname, *port, *database))
	if err != nil {
		log.Fatal(err)
	}

	// GET /book/1
	app.Get("/book/:id", func(c *fiber.Ctx) error {
		book := &Book{}
		bookId := c.Params("id")
		bookRow := dbPool.QueryRow(context.Background(), "SELECT * FROM books WHERE id=$1;", bookId)
		err := bookRow.Scan(&book.Id, &book.Title, &book.Author)
		if err != nil {
			return err
		}

		return c.JSON(book)
	})

	// GET /books
	app.Get("/books", func(c *fiber.Ctx) error {
		books := []*Book{}
		booksRows, err := dbPool.Query(context.Background(), "SELECT * FROM books;")
		if err != nil {
			return err
		}

		for booksRows.Next() {
			book := &Book{}
			err = booksRows.Scan(&book.Id, &book.Title, &book.Author)
			if err != nil {
				return err
			}

			books = append(books, book)
		}

		return c.JSON(books)
	})

	// POST /book/create
	app.Post("/book/create", func(c *fiber.Ctx) error {
		book := &Book{}
		json.NewDecoder(bytes.NewReader(c.Body())).Decode(&book)
		_, err := dbPool.Exec(context.Background(), `INSERT INTO books(title, author) VALUES ($1, $2);`, &book.Title, &book.Author)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusOK)
	})

	// PUT /book/update
	app.Put("/book/update", func(c *fiber.Ctx) error {
		book := &Book{}
		json.NewDecoder(bytes.NewReader(c.Body())).Decode(&book)
		_, err := dbPool.Exec(context.Background(), `UPDATE books SET title=$1, author=$2 WHERE id=$3;`, book.Title, book.Author, book.Id)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusOK)
	})

	// DELETE /book/delete/1
	app.Delete("/book/delete/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, err := dbPool.Exec(context.Background(), "DELETE FROM books WHERE id=$1", id)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusOK)
	})

	log.Fatal(app.Listen(":8080"))
}
