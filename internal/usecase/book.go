package usecase

import (
	"github.com/cothromachd/books-api/internal/entity"
)

type Storage interface {
	GetBook(string) (entity.Book, error)
	GetBooks() ([]entity.Book, error)
	PostBook(title, author string) (string, error)
	UpdateBook(id string, book entity.Book) error
	DeleteBook(id string) error
}

type Cache interface {
	HasBook(id string) (bool, error)
	GetBook(id string) (entity.Book, error)
	SetBook(id string, book entity.Book) error
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
	exists, err := uc.cache.HasBook(id)
	if err != nil {
		return entity.Book{}, err
	}
	if exists {
		return uc.cache.GetBook(id)
	}
	
	book, err := uc.storage.GetBook(id)
	if err != nil {
		return entity.Book{}, err
	}

	return book, uc.cache.SetBook(id, book)
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

func (uc *Book) UpdateBook(id string, book entity.Book) error  {
	err := uc.storage.UpdateBook(id, book)
	if err != nil {
		return err
	}

	err = uc.cache.SetBook(id, book)
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