package app

import (
	"temp/internal/app/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type App struct{}

func New() App {
	return App{}
}

func (a *App) Run() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		r := Response{
			Result: models.Tmp{
				ID:      uuid.New().String(),
				Message: "Faggot",
			},
		}
		return c.JSON(r)
	})

	app.Listen(":8000")
}
