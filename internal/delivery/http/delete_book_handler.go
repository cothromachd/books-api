package delivery

import "github.com/gofiber/fiber/v2"

func (h *Handler) DeleteBook(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	err := h.uc.DeleteBook(id)
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}