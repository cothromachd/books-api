package app

import (
	"log"

	"github.com/cothromachd/books-api/internal/config"
	"github.com/cothromachd/books-api/internal/delivery/http"
	"github.com/cothromachd/books-api/internal/infrastructure/repository"
	"github.com/cothromachd/books-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

const (
	configPath = "/Users/khalidmagnificent/Desktop/goground/books-api/configs/config.yaml"
)

func Run() error {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		},
	})

	cfg, err := config.New(configPath)
	if err != nil {
		return err
	}

	s, err := repo.NewPostgresStorage(cfg)
	if err != nil {
		return err
	}
	defer s.Close()

	c := repo.NewRedisCache(cfg)

	uc := usecase.NewBook(s, c)

	app = delivery.NewHandler(app, uc)

	return(app.Listen(cfg.API.Host))
}