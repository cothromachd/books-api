package delivery

import (
	"github.com/cothromachd/books-api/internal/entity"
	"github.com/gofiber/fiber/v2"
)

type BookUseCase interface {
	GetBooks() ([]entity.Book, error)
	GetBook(id string) (entity.Book, error)
	PostBook(book entity.Book) error
	UpdateBook(id string, book entity.Book) error
	DeleteBook(id string) error
}

type Handler struct {
	uc BookUseCase
}

func NewHandler(app *fiber.App, uc BookUseCase) *fiber.App {
	h := Handler{
		uc: uc,
	}

	router := app.Group("book")

	router.Get("/all", h.GetBooks)
	router.Get("/:id", h.GetBook)
	router.Post("/create", h.CreateBook)
	router.Put("/update/:id", h.UpdateBook)
	router.Delete("/delete/:id", h.DeleteBook)

	return app
}
