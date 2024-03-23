package middleware

import (
	"errors"
	"strings"
	manager "temp/internal/app/domain/managers"
	model "temp/internal/app/domain/models"
	service "temp/internal/app/domain/services"
	constant "temp/internal/app/router/constants"
	controller "temp/internal/app/router/controllers"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrAuthorizationHeaderRequired = fiber.NewError(fiber.StatusBadRequest, "Authorization header is required")
	ErrInvalidAuthorizationHeader  = fiber.NewError(fiber.StatusBadRequest, "invalid Authorization header")
	ErrUserAgentNotPassed          = fiber.NewError(fiber.StatusBadRequest, "User-Agent header is required")
	ErrFingerprintNotPassed        = fiber.NewError(fiber.StatusBadRequest, "X-Fingerprint header is required")
	ErrSessionNotFound             = fiber.NewError(fiber.StatusUnauthorized, "session assosiated with token not found, probably deleted")
	ErrUnknownAccessToken          = fiber.NewError(fiber.StatusNotFound, "unknown access token was used, your session dropped")
	ErrUnknownDevice               = fiber.NewError(fiber.StatusNotFound, "not associated with refresh token device used, your session dropped")
)

type tokenValidatorFn func(t string) (model.TokenPayload, error)

type AuthGuard struct {
	s  service.SessionsService
	tm manager.TokenManager
}

func NewAuthGuard(s service.SessionsService, tm manager.TokenManager) *AuthGuard {
	return &AuthGuard{s, tm}
}

func (a *AuthGuard) RefreshGuard(c *fiber.Ctx) error {
	return a.validate(c, a.tm.Parse)
}

func (a *AuthGuard) AccessGuard(c *fiber.Ctx) error {
	return a.validate(c, a.tm.ParseAndValidate)
}

func (a *AuthGuard) validate(c *fiber.Ctx, fn tokenValidatorFn) error {
	accessToken, err := a.retrieveToken(c)
	if err != nil {
		return err
	}

	tokenPayload, err := fn(accessToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	uagent := c.Get(fiber.HeaderUserAgent)
	if uagent == "" {
		return ErrUserAgentNotPassed
	}

	fprint := c.Get(constant.HeaderFingerprint)
	if fprint == "" {
		return ErrFingerprintNotPassed
	}

	err = a.s.Validate(tokenPayload.SessionID, accessToken, uagent, fprint)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return controller.ErrInvalidToken
		}

		return controller.ErrUnauthorized(err)
	}

	c.Locals(constant.LocalKeyCommon, model.CommonRequestPayload{
		TokenPayload: tokenPayload,
		DeviceInfo: model.DeviceInfo{
			Uagent: uagent,
			Fprint: fprint,
		},
	})

	return c.Next()
}

func (a *AuthGuard) retrieveToken(c *fiber.Ctx) (string, error) {
	h := c.Get(fiber.HeaderAuthorization)
	if h == "" {
		return "", ErrAuthorizationHeaderRequired
	}

	parts := strings.Split(h, " ")
	if len(parts) != 2 {
		return "", ErrInvalidAuthorizationHeader
	}
	if parts[0] != "Bearer" {
		return "", ErrInvalidAuthorizationHeader
	}
	return parts[1], nil
}
