package controller

import (
	"errors"
	service "temp/internal/app/domain/services"
	"temp/internal/app/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailNotValid    = fiber.NewError(fiber.StatusBadRequest, "email is not valid")
	ErrPasswordTooSmall = fiber.NewError(fiber.StatusBadRequest, "password cannot be shorter 8 symbols")
	ErrPasswordTooLong  = fiber.NewError(fiber.StatusBadRequest, "password cannot be longer 72 symbols")

	ErrUserIDNotPassed = fiber.NewError(fiber.StatusNotFound, "user_id cannot be empty")
	ErrUserNotFound    = fiber.NewError(fiber.StatusNotFound, "user not found")

	ErrUserExists = fiber.NewError(fiber.StatusConflict, "user already exists")
)

type Users struct {
	s service.UsersService
}

func NewUsers(s service.UsersService) *Users {
	return &Users{s}
}

func (u *Users) GetInfo(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return ErrUserIDNotPassed
	}

	user, err := u.s.FindDetailedByID(userID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return ErrUserNotFound
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(user)
}

type createRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type createResponse struct {
	ID string `json:"user_id"`
}

func (u *Users) Create(c *fiber.Ctx) error {
	var body createRequest
	if err := c.BodyParser(&body); err != nil {
		return ErrInvalidBody(err)
	}

	if !utils.IsEmailValid(body.Email) {
		return ErrEmailNotValid
	}

	if len(body.Password) < 8 {
		return ErrPasswordTooSmall
	}

	userID, err := u.s.Create(body.Email, body.Password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return ErrPasswordTooLong
		}

		if errors.Is(err, service.ErrAlreadyExists) {
			return ErrUserExists
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(createResponse{userID})
}

func (u *Users) Update(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

func (u *Users) Delete(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

func (u *Users) ChangeEmail(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}
