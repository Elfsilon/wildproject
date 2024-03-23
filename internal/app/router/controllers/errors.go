package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func ErrInvalidBody(err error) error {
	message := fmt.Sprintf("invalid body: %s", err)
	return fiber.NewError(fiber.StatusBadRequest, message)
}

func ErrUnauthorized(err error) error {
	message := fmt.Sprintf("unauthorized: %s", err)
	return fiber.NewError(fiber.StatusUnauthorized, message)
}
