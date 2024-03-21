package ctr

import (
	"database/sql"
	"errors"
	"net/mail"
	rep "temp/internal/app/repositories"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserIDNotPassed = fiber.NewError(fiber.StatusNotFound, "user_id cannot be empty")
	ErrUserNotFound    = fiber.NewError(fiber.StatusNotFound, "user not found")
	ErrEmailNotValid   = fiber.NewError(fiber.StatusBadRequest, "email is not valid")
	ErrUserExists      = fiber.NewError(fiber.StatusConflict, "user already exists")
	ErrPasswordTooLong = fiber.NewError(fiber.StatusBadRequest, "password cannot be longer 72 symbols")
)

type Users struct {
	r *rep.Users
}

func NewUsers(r *rep.Users) *Users {
	return &Users{r}
}

func (u *Users) isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (u *Users) GetInfo(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return ErrUserIDNotPassed
	}

	detailed, err := u.r.GetDetailed(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(detailed)
}

func (u *Users) Create(c *fiber.Ctx) error {
	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !u.isEmailValid(payload.Email) {
		return ErrEmailNotValid
	}

	phash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 16)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return ErrPasswordTooLong
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	userID, err := u.r.Create(payload.Email, string(phash))
	if err != nil {
		if errors.Is(err, rep.ErrUserExists) {
			return ErrUserExists
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	res := struct {
		ID string `json:"user_id"`
	}{userID}

	return c.JSON(res)
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
