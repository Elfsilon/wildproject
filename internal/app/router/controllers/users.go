package controller

import (
	"errors"
	model "wildproject/internal/app/domain/models"
	service "wildproject/internal/app/domain/services"
	constant "wildproject/internal/app/router/constants"
	"wildproject/internal/app/utils"

	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailNotValid    = fiber.NewError(fiber.StatusBadRequest, "email is not valid")
	ErrPasswordTooSmall = fiber.NewError(fiber.StatusBadRequest, "password cannot be shorter 8 symbols")
	ErrPasswordTooLong  = fiber.NewError(fiber.StatusBadRequest, "password cannot be longer 72 symbols")

	ErrUserNotFound = fiber.NewError(fiber.StatusNotFound, "user not found")

	ErrUserExists = fiber.NewError(fiber.StatusConflict, "user already exists")
)

type Users struct {
	s service.UsersService
}

func NewUsers(s service.UsersService) *Users {
	return &Users{s}
}

func (u *Users) GetInfo(c *fiber.Ctx) error {
	hub := fibersentry.GetHubFromContext(c)

	p, ok := c.Locals(constant.LocalKeyCommon).(model.CommonRequestPayload)
	if !ok {
		return ErrInvalidCommonPayload
	}

	user, err := u.s.FindDetailedByID(p.UserID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return ErrUserNotFound
		}

		hub.CaptureException(err)
		return err
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
	hub := fibersentry.GetHubFromContext(c)

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

		hub.CaptureException(err)
		return err
	}

	return c.JSON(createResponse{userID})
}

type changeNameRequest struct {
	Name string `json:"name"`
}

func (u *Users) ChangeName(c *fiber.Ctx) error {
	var request changeNameRequest

	if err := c.BodyParser(&request); err != nil {
		return ErrInvalidBody(err)
	}

	p, ok := c.Locals(constant.LocalKeyCommon).(model.CommonRequestPayload)
	if !ok {
		return ErrInvalidCommonPayload
	}

	name, err := u.s.ChangeName(p.UserID, request.Name)
	if err != nil {
		return err
	}

	return c.JSON(changeNameRequest{
		Name: name,
	})
}

type changeSexRequest struct {
	SexID int `json:"sex_id"`
}

func (u *Users) ChangeSex(c *fiber.Ctx) error {
	var request changeSexRequest

	if err := c.BodyParser(&request); err != nil {
		return ErrInvalidBody(err)
	}

	p, ok := c.Locals(constant.LocalKeyCommon).(model.CommonRequestPayload)
	if !ok {
		return ErrInvalidCommonPayload
	}

	sexID, err := u.s.ChangeSex(p.UserID, request.SexID)
	if err != nil {
		return err
	}

	return c.JSON(changeSexRequest{
		SexID: sexID,
	})
}

type changeEmailRequest struct {
	Email string `json:"email"`
}

func (u *Users) ChangeEmail(c *fiber.Ctx) error {
	var request changeEmailRequest

	if err := c.BodyParser(&request); err != nil {
		return ErrInvalidBody(err)
	}

	if !utils.IsEmailValid(request.Email) {
		return ErrEmailNotValid
	}

	p, ok := c.Locals(constant.LocalKeyCommon).(model.CommonRequestPayload)
	if !ok {
		return ErrInvalidCommonPayload
	}

	email, err := u.s.ChangeEmail(p.UserID, request.Email)
	if err != nil {
		return err
	}

	return c.JSON(changeEmailRequest{
		Email: email,
	})
}

type changePasswordRequest struct {
	Password string `json:"password"`
}

func (u *Users) ChangePassword(c *fiber.Ctx) error {
	var request changePasswordRequest

	if err := c.BodyParser(&request); err != nil {
		return ErrInvalidBody(err)
	}

	if len(request.Password) < 8 {
		return ErrPasswordTooSmall
	}

	p, ok := c.Locals(constant.LocalKeyCommon).(model.CommonRequestPayload)
	if !ok {
		return ErrInvalidCommonPayload
	}

	err := u.s.ChangePassword(p.UserID, request.Password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return ErrPasswordTooLong
		}

		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (u *Users) ChangeImage(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

func (u *Users) Delete(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}
