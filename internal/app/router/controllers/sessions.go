package controller

import (
	"errors"
	"strconv"
	model "wildproject/internal/app/domain/models"
	service "wildproject/internal/app/domain/services"
	constant "wildproject/internal/app/router/constants"

	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrInvalidSessionID     = fiber.NewError(fiber.StatusBadRequest, "invalid session_id")
	ErrInvalidDevice        = fiber.NewError(fiber.StatusBadRequest, "invalid locals device_info")
	ErrWrongEmailOrPassword = fiber.NewError(fiber.StatusBadRequest, "wrong email or password")
	ErrUserAgentNotPassed   = fiber.NewError(fiber.StatusBadRequest, "User-Agent header is required")
	ErrFingerprintNotPassed = fiber.NewError(fiber.StatusBadRequest, "X-Fingerprint header is required")

	ErrAlreadyAuthorized = fiber.NewError(fiber.StatusConflict, "passed token is valid and so user is already authorized")

	ErrUnknownRefreshToken = fiber.NewError(fiber.StatusNotFound, "unknown refresh token was used, your session dropped")
	ErrExpiredRefreshToken = fiber.NewError(fiber.StatusNotFound, "refresh token is expired, your session dropped")
	ErrSessionNotFound     = fiber.NewError(fiber.StatusNotFound, "session not found")

	ErrInvalidToken = fiber.NewError(fiber.StatusUnauthorized, "invalid token")
)

type Sessions struct {
	sSer service.SessionsService
	uSer service.UsersService
}

func NewSessions(
	sSer service.SessionsService,
	uSer service.UsersService,
) *Sessions {
	return &Sessions{sSer, uSer}
}

func (s *Sessions) GetByID(c *fiber.Ctx) error {
	hub := fibersentry.GetHubFromContext(c)

	sessionID, err := strconv.Atoi(c.Params("session_id"))
	if err != nil {
		return ErrInvalidSessionID
	}

	session, err := s.sSer.Find(sessionID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return ErrSessionNotFound
		}

		hub.CaptureException(err)
		return err
	}

	return c.JSON(session)
}

type sessionsResponse struct {
	Sessions []model.ClientRefreshSession `json:"sessions"`
}

// Get all active user's refresh sessions
func (s *Sessions) GetAllByUserID(c *fiber.Ctx) error {
	hub := fibersentry.GetHubFromContext(c)

	p, ok := c.Locals(constant.LocalKeyCommon).(model.CommonRequestPayload)
	if !ok {
		return ErrInvalidCommonPayload
	}

	sessions, err := s.sSer.FindAll(p.UserID, "", "")
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return ErrSessionNotFound
		}

		hub.CaptureException(err)
		return err
	}

	return c.JSON(sessionsResponse{sessions})
}

type newSessionRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Creates new refresh session for the user if valid Authorization header is not passed
func (s *Sessions) Create(c *fiber.Ctx) error {
	hub := fibersentry.GetHubFromContext(c)

	// TODO: Check if user already authorized
	// h := c.Get(fiber.HeaderAuthorization)
	// hIsValid := false
	// if h != "" && hIsValid {
	// 	return ErrAlreadyAuthorized
	// }

	var request newSessionRequest

	if err := c.BodyParser(&request); err != nil {
		return ErrInvalidBody(err)
	}

	// TODO: add email and password validation
	// minlength: 8;
	// maxlength: 72;
	// required: lower;
	// required: upper;
	// required: digit;
	// required: [#$%&*.@^];

	userID, err := s.uSer.Authenticate(request.Email, request.Password)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) || errors.Is(err, service.ErrPasswordsMismatch) {
			return ErrWrongEmailOrPassword
		}

		hub.CaptureException(err)
		return err
	}

	// Check if new device is used else drop an old session
	uagent := c.Get(fiber.HeaderUserAgent)
	if uagent == "" {
		return ErrUserAgentNotPassed
	}

	fprint := c.Get("X-Fingerprint")
	if fprint == "" {
		return ErrFingerprintNotPassed
	}

	tokens, err := s.sSer.Create(userID, uagent, fprint)
	if err != nil {
		hub.CaptureException(err)
		return err
	}

	return c.JSON(tokens)
}

type resfreshSessionRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Generates a new pair of resfresh + access tokens by valid refresh token
func (s *Sessions) Refresh(c *fiber.Ctx) error {
	hub := fibersentry.GetHubFromContext(c)

	payload, ok := c.Locals(constant.LocalKeyCommon).(model.CommonRequestPayload)
	if !ok {
		return ErrInvalidCommonPayload
	}
	userID := payload.UserID

	var request resfreshSessionRequest

	if err := c.BodyParser(&request); err != nil {
		return ErrInvalidBody(err)
	}

	cp, ok := c.Locals(constant.LocalKeyCommon).(model.CommonRequestPayload)
	if !ok {
		return ErrInvalidDevice
	}

	tokens, err := s.sSer.Refresh(
		request.RefreshToken,
		userID,
		cp.Uagent,
		cp.Fprint,
	)

	if err != nil {
		if errors.Is(err, service.ErrUnknownToken) {
			return ErrUnknownRefreshToken
		}

		if errors.Is(err, service.ErrExpiredToken) {
			return ErrExpiredRefreshToken
		}

		hub.CaptureException(err)
		return err
	}

	return c.JSON(tokens)
}

// Drops specified user's refresh session
func (s *Sessions) Drop(c *fiber.Ctx) error {
	hub := fibersentry.GetHubFromContext(c)

	sessionID, err := strconv.Atoi(c.Params("session_id"))
	if err != nil {
		return ErrInvalidSessionID
	}

	if err := s.sSer.Drop(sessionID); err != nil {
		hub.CaptureException(err)
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

// Drops all user's refresh sessions (equivalent to logout)
func (s *Sessions) DropAll(c *fiber.Ctx) error {
	hub := fibersentry.GetHubFromContext(c)

	p, ok := c.Locals(constant.LocalKeyCommon).(model.CommonRequestPayload)
	if !ok {
		return ErrInvalidCommonPayload
	}

	if err := s.sSer.DropAll(p.UserID, "", ""); err != nil {
		hub.CaptureException(err)
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
