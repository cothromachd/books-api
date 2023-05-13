package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Book struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var dsn = flag.String("dsn", "postgres://username:password@hostname:port/database-name", "Connect to database string\nExample: postgres://{username}:{password}@{hostname}:{port}/{database-name}")
var rdbAddr = flag.String("ra", "redis-addres:redis-port", "Redis addres\n{ip}:{port}")

func main() {
	flag.Parse()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		},
	})

	rdb := redis.NewClient(&redis.Options{
		Addr:     *rdbAddr,
		Password: "",
		DB:       0,
	})

	dbPool, err := pgxpool.New(context.Background(), *dsn)
	if err != nil {
		log.Fatal(err)
	}

	// GET /book/1
	app.Get("/book/:id", func(c *fiber.Ctx) error {
		book := &Book{}
		bookId := c.Params("id")
		bookJson, err := rdb.Get(context.Background(), bookId).Result()
		if err == redis.Nil {
			bookRow := dbPool.QueryRow(context.Background(), "SELECT * FROM books WHERE id=$1;", bookId)
			err := bookRow.Scan(&book.Id, &book.Title, &book.Author)
			if err != nil {
				return err
			}

			bookJson, err := json.Marshal(book)
			if err != nil {
				return err
			}

			rdb.Set(context.Background(), bookId, string(bookJson), time.Hour)
		} else if err == nil {
			return c.JSON(bookJson)
		} else {
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
		err := json.NewDecoder(bytes.NewReader(c.Body())).Decode(&book)
		if err != nil {
			return err
		}

		idRow := dbPool.QueryRow(context.Background(), `INSERT INTO books(title, author) VALUES ($1, $2) RETURNING id;`, &book.Title, &book.Author)

		var id int
		err = idRow.Scan(&id)
		if err != nil {
			return err
		}

		rdb.Set(context.Background(), strconv.Itoa(id), string(c.Body()), time.Hour)
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
		
		err := rdb.Del(context.Background(), "id").Err()
		if err != nil {
			return err
		}

		_, err = dbPool.Exec(context.Background(), "DELETE FROM books WHERE id=$1", id)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusOK)
	})

	log.Fatal(app.Listen(":8080"))
}
