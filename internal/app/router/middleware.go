package router

import (
	"strings"
	tokenmanager "temp/internal/app/token-manager"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrAuthorizationHeaderRequired = fiber.NewError(fiber.StatusBadRequest, "Authorization header is required")
	ErrInvalidAuthorizationHeader  = fiber.NewError(fiber.StatusBadRequest, "invalid Authorization header")
)

type AuthGuard struct {
	tm *tokenmanager.TokenManager
}

func NewAuthGuard(tm *tokenmanager.TokenManager) *AuthGuard {
	return &AuthGuard{tm}
}

func (a *AuthGuard) ValidateToken(c *fiber.Ctx) error {
	h := c.Get(fiber.HeaderAuthorization)
	if h == "" {
		return ErrAuthorizationHeaderRequired
	}

	parts := strings.Split(h, " ")
	if len(parts) != 2 {
		return ErrInvalidAuthorizationHeader
	}
	if parts[0] != "Bearer" {
		return ErrInvalidAuthorizationHeader
	}
	accessToken := parts[1]

	sessionID, userID, err := a.tm.Validate(accessToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	c.Locals("session_id", sessionID)
	c.Locals("user_id", userID)

	return c.Next()
}
