package controller

import (
	model "temp/internal/app/domain/models"

	"github.com/gofiber/fiber/v2"
)

func HealthCheck(c *fiber.Ctx) error {
	r := model.Response{
		Result: model.Status{
			Message: "Hi mafaka!",
		},
	}
	return c.JSON(r)
}
