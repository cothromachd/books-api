package delivery

import (
	"github.com/cothromachd/books-api/internal/entity"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateBook(ctx *fiber.Ctx) error {
	bookJson := ctx.Body()
	book, err := entity.Unmap(string(bookJson))
	if err != nil {
		return err
	}

	err = h.uc.PostBook(book)
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}