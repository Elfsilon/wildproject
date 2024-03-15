package ctr

import (
	m "temp/internal/app/models"

	"github.com/gofiber/fiber/v2"
)

func HealthCheck(c *fiber.Ctx) error {
	r := m.Response{
		Result: m.Status{
			Message: "Hi mafaka!",
		},
	}
	return c.JSON(r)
}
