package usecase

import (
	"github.com/cothromachd/books-api/internal/entity"
)

type Storage interface {
	GetBook(string) (entity.Book, error)
	GetBooks() ([]entity.Book, error)
	PostBook(title, author string) (int, error)
	UpdateBook(book entity.Book) error
	DeleteBook(id string) error
}

type Cache interface {
	GetBook(id string) (entity.Book, error)
	SetBook(id int, book entity.Book) error
	DeleteBook(string) error
}

type Book struct {
	storage Storage
	cache Cache
}

func NewBook(storage Storage, cache Cache) *Book {
	return &Book{
		storage: storage,
		cache: cache,
	}
}

func (uc *Book) GetBooks() ([]entity.Book, error) {
	books, err := uc.storage.GetBooks()
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (uc *Book) GetBook(id string) (entity.Book, error) {
	book, err := uc.cache.GetBook(id)
	if err != nil && err.Error() == "redis: nil" {
		book, err := uc.storage.GetBook(id)
		if err != nil {
			return entity.Book{}, err
		}

		return book, nil
	} else if err != nil {
		return entity.Book{}, err
	} else {
		return book, nil
	}
}

func (uc *Book) PostBook(book entity.Book) error {
	id, err := uc.storage.PostBook(book.Title, book.Author)
	if err != nil {
		return err
	}

	err = uc.cache.SetBook(id, book)
	if err != nil {
		return err
	}

	return nil
}

func (uc *Book) UpdateBook(book entity.Book) error  {
	err := uc.storage.UpdateBook(book)
	if err != nil {
		return err
	}

	err = uc.cache.SetBook(book.Id, book)
	if err != nil {
		return err
	}

	return nil
}

func (uc *Book) DeleteBook(id string) error  {
	err := uc.storage.DeleteBook(id)
	if err != nil {
		return err
	}

	err = uc.cache.DeleteBook(id)
	if err != nil {
		return err
	}

	return nil
}