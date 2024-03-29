package middleware

import (
	"errors"
	"fmt"
	"strings"
	manager "wildproject/internal/app/domain/managers"
	model "wildproject/internal/app/domain/models"
	service "wildproject/internal/app/domain/services"
	constant "wildproject/internal/app/router/constants"
	controller "wildproject/internal/app/router/controllers"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/contrib/fibersentry"
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
	hub := fibersentry.GetHubFromContext(c)

	accessToken, err := a.retrieveToken(c)
	if err != nil {
		return err
	}

	tokenPayload, err := fn(accessToken)
	if err != nil {
		hub.CaptureException(err)
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

		hub.CaptureException(err)
		return controller.ErrUnauthorized(err)
	}

	c.Locals(constant.LocalKeyCommon, model.CommonRequestPayload{
		TokenPayload: tokenPayload,
		DeviceInfo: model.DeviceInfo{
			Uagent: uagent,
			Fprint: fprint,
		},
	})

	hub.Scope().SetUser(sentry.User{
		ID: tokenPayload.UserID,
		Data: map[string]string{
			"session_id": fmt.Sprint(tokenPayload.SessionID),
		},
	})

	hub.Scope().SetTags(map[string]string{
		"User-Agent": uagent,
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
