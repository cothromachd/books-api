package entity

import (
	"bytes"
	"encoding/json"
)

type Book struct {
	Id     int    `json:"id,omitempty"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// func ValidateUser(b *Book) bool { return true }

func (book *Book) Map() (string, error) {
	bookJson, err := json.Marshal(book)
	if err != nil {
		return "", err
	}

	return string(bookJson), nil
}

func Unmap(bookJson string) (Book, error) {
	var book Book
	err := json.NewDecoder(bytes.NewReader([]byte(bookJson))).Decode(&book)
	if err != nil {
		return Book{}, err
	}
	return book, nil
}