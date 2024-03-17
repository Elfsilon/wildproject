package ctr

import (
	rep "temp/internal/app/repositories"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	r *rep.Users
}

func NewUsers(r *rep.Users) *Users {
	return &Users{r}
}

func (u *Users) GetInfo(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "user_id cannot be empty")
	}

	detailed, err := u.r.GetDetailed(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(detailed)
}

// POST users/ { email, password }
func (u *Users) Create(c *fiber.Ctx) error {
	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	phash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 16)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	userID, err := u.r.Create(payload.Email, string(phash))
	if err != nil {
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
