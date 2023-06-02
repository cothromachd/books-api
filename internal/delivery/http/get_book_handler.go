package delivery

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetBook(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	book, err := h.uc.GetBook(id)
	if err != nil {
		return err
	}

	return ctx.JSON(book)
}