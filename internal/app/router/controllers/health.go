package controller

import (
	"github.com/gofiber/fiber/v2"
)

type healthResponse struct {
	Message string `json:"message"`
}

func HealthCheck(c *fiber.Ctx) error {
	r := healthResponse{
		Message: "Hi mafaka!",
	}
	return c.JSON(r)
}
