package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrInvalidCommonPayload = fiber.NewError(fiber.StatusInternalServerError, "invalid standart request payload")
	ErrUserIDNotPassed      = fiber.NewError(fiber.StatusNotFound, "user_id cannot be empty")
)

func ErrInvalidBody(err error) error {
	message := fmt.Sprintf("invalid body: %s", err)
	return fiber.NewError(fiber.StatusBadRequest, message)
}

func ErrUnauthorized(err error) error {
	message := fmt.Sprintf("unauthorized: %s", err)
	return fiber.NewError(fiber.StatusUnauthorized, message)
}
