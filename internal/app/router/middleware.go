package router

import (
	"database/sql"
	"errors"
	"strings"
	model "temp/internal/app/domain/models"
	rep "temp/internal/app/repositories"
	tokenmanager "temp/internal/app/token-manager"

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

type AuthGuard struct {
	s  *rep.Sessions
	tm *tokenmanager.TokenManager
}

func NewAuthGuard(tm *tokenmanager.TokenManager, s *rep.Sessions) *AuthGuard {
	return &AuthGuard{s, tm}
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

func (a *AuthGuard) validateSessionAndDevice(
	c *fiber.Ctx,
	accessToken string,
	sessionID int,
) error {
	uagent := c.Get(fiber.HeaderUserAgent)
	if uagent == "" {
		return ErrUserAgentNotPassed
	}

	fprint := c.Get("X-Fingerprint")
	if fprint == "" {
		return ErrFingerprintNotPassed
	}

	oldSession, err := a.s.FindBySessionID(sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrSessionNotFound
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if accessToken != oldSession.AccessToken {
		a.s.Drop(oldSession.SessionID)
		return ErrUnknownAccessToken
	}

	if uagent != oldSession.Uagent || fprint != oldSession.Fprint {
		a.s.Drop(oldSession.SessionID)
		return ErrUnknownDevice
	}

	c.Locals("session", oldSession)
	c.Locals("device_info", model.DeviceInfo{
		Uagent: uagent,
		Fprint: fprint,
	})

	return nil
}

func (a *AuthGuard) RefreshGuard(c *fiber.Ctx) error {
	return a.validate(c, a.tm.GetClaims)
}

func (a *AuthGuard) AccessGuard(c *fiber.Ctx) error {
	return a.validate(c, a.tm.Validate)
}

func (a *AuthGuard) validate(
	c *fiber.Ctx,
	meth func(t string) (model.TokenData, error),
) error {
	accessToken, err := a.retrieveToken(c)
	if err != nil {
		return err
	}

	tokenData, err := meth(accessToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	err = a.validateSessionAndDevice(c, accessToken, tokenData.SessionID)
	if err != nil {
		return err
	}

	c.Locals("access_token", accessToken)
	c.Locals("token_data", tokenData)

	return c.Next()
}
