package delivery

import (
	"github.com/cothromachd/books-api/internal/entity"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) UpdateBook(ctx *fiber.Ctx) error {
	bookJson := ctx.Body()
	book, err := entity.Unmap(string(bookJson))
	if err != nil {
		return err
	}

	err = h.uc.UpdateBook(book)
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}