package repo

import (
	"context"
	"strconv"

	"github.com/cothromachd/books-api/internal/config"
	"github.com/cothromachd/books-api/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pgPool *pgxpool.Pool
}

func NewPostgresStorage(cfg *config.Config) (*PostgresStorage, error) {
	pgPool, err := pgxpool.New(context.Background(), cfg.DB.Conn)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{pgPool: pgPool}, nil
}

func (s PostgresStorage) Close() {
	s.pgPool.Close()
}

func (s *PostgresStorage) GetBook(id string) (entity.Book, error) {
	book := entity.Book{}
	bookRow := s.pgPool.QueryRow(context.Background(), "SELECT * FROM books WHERE id=$1;", id)
	err := bookRow.Scan(&book.Id, &book.Title, &book.Author)
	if err != nil {
		return entity.Book{}, err
	}

	return book, nil
}

func (s *PostgresStorage) GetBooks() ([]entity.Book, error) {
	books := []entity.Book{}
	booksRows, err := s.pgPool.Query(context.Background(), "SELECT * FROM books;")
	if err != nil {
		return nil, err
	}

	for booksRows.Next() {
		book := entity.Book{}
		err = booksRows.Scan(&book.Id, &book.Title, &book.Author)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

func (s *PostgresStorage) PostBook(title, author string) (string, error) {
	idRow := s.pgPool.QueryRow(context.Background(), `INSERT INTO books(title, author) VALUES ($1, $2) RETURNING id;`, title, author)

	var bookId int
	err := idRow.Scan(&bookId)
	if err != nil {
		return "", err
	}

	id := strconv.Itoa(bookId)
	return id, nil
}

func (s *PostgresStorage) UpdateBook(id string, book entity.Book) error {
	_, err := s.pgPool.Exec(context.Background(), `UPDATE books SET title=$1, author=$2 WHERE id=$3;`, book.Title, book.Author, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) DeleteBook(id string) error {
	_, err := s.pgPool.Exec(context.Background(), "DELETE FROM books WHERE id=$1", id)
	if err != nil {
		return err
	}

	return nil
}