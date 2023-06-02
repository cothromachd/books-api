package delivery

import "github.com/gofiber/fiber/v2"

func (h *Handler) GetBooks(ctx *fiber.Ctx) error {
	books, err := h.uc.GetBooks()
	if err != nil {
		return err
	}

	return ctx.JSON(books)
}