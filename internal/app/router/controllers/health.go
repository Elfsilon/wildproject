package controller

import (
	"github.com/gofiber/fiber/v2"
)

type healthResponse struct {
	Message string `json:"message"`
}

func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(healthResponse{
		Message: "Hi mafaka!",
	})
}
