package ctr

import (
	"database/sql"
	"errors"
	"fmt"
	m "temp/internal/app/models"
	rep "temp/internal/app/repositories"
	tokenmanager "temp/internal/app/token-manager"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrWrongEmailOrPassword = fiber.NewError(fiber.StatusInternalServerError, "wrong email or password")
	ErrMismatchCurrentUser  = fiber.NewError(fiber.StatusInternalServerError, "current_user mismatches type of User")
	ErrAlreadyAuthorized    = fiber.NewError(fiber.StatusConflict, "passed token is valid and so user is already authorized")
	ErrUserAgentNotPassed   = fiber.NewError(fiber.StatusBadRequest, "User-Agent header is required")
	ErrFingerprintNotPassed = fiber.NewError(fiber.StatusBadRequest, "X-Fingerprint header is required")
)

type Sessions struct {
	cfg *m.AuthConfig
	r   *rep.Sessions
	ur  *rep.Users
	tm  *tokenmanager.TokenManager
}

func NewSessions(
	cfg *m.AuthConfig,
	r *rep.Sessions,
	ur *rep.Users,
	tm *tokenmanager.TokenManager,
) *Sessions {
	return &Sessions{cfg, r, ur, tm}
}

// Get all active user's refresh sessions
func (s *Sessions) GetAllByUserID(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return ErrUserIDNotPassed
	}

	sessions, err := s.r.FindByUserID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(struct {
		Sessions []m.RefreshSession `json:"sessions"`
	}{
		Sessions: sessions,
	})
}

// Creates new refresh session for the user if valid Authorization header is not passed
func (s *Sessions) New(c *fiber.Ctx) error {
	// TODO: Check if user already authorized
	// h := c.Get(fiber.HeaderAuthorization)
	// hIsValid := false
	// if h != "" && hIsValid {
	// 	return ErrAlreadyAuthorized
	// }

	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&payload); err != nil {
		message := fmt.Sprintf("invalid body: %s", err)
		return fiber.NewError(fiber.StatusBadRequest, message)
	}

	// TODO: add email and password validation
	// minlength: 8;
	// maxlength: 72;
	// required: lower;
	// required: upper;
	// required: digit;
	// required: [#$%&*.@^];

	// Authorize user
	userID, cred, err := s.ur.FindCredentialsByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrWrongEmailOrPassword
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bHash := []byte(cred.PasswordHash)
	bPass := []byte(payload.Password)

	if err := bcrypt.CompareHashAndPassword(bHash, bPass); err != nil {
		return ErrWrongEmailOrPassword
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

	sessionIDs, err := s.r.FindByDevice(userID, uagent, fprint)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	log.Infof("dropping sessions: %v", sessionIDs)
	if sessionIDs != nil && len(sessionIDs) > 0 {
		// Try to delete sessions with this device
		for _, sessionID := range sessionIDs {
			if err := s.r.Drop(sessionID); err != nil {
				// TODO: Figure out what to do
				continue
			}
		}
	} else {
		// TODO: push notification: New device login detected
	}

	// Generate new refresh session
	rTokenExriresAt := time.Now().Add(s.cfg.RefreshTokenTTL).UTC()
	sessionID, refreshToken, err := s.r.Create(userID, uagent, fprint, rTokenExriresAt)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	accessToken, err := s.tm.Generate(sessionID, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = s.r.SetAccessToken(sessionID, accessToken)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(m.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Generates a new pair of resfresh + access tokens by valid refresh token
func (s *Sessions) Refresh(c *fiber.Ctx) error {
	// Need: access_token, refresh_token, user_agent, fingerprint
	//
	// Checks that old access token matches the one from saved session
	//	 If not: drops session and sends "unknown access token was used, your session dropped"
	// Checks that old device is used
	//   If not: drops session and sends "not associated with refresh token device used, your session dropped"
	// Removes old session
	// Created new session

	return fiber.ErrNotImplemented
}

// Drops specidied user's refresh session
func (s *Sessions) Drop(c *fiber.Ctx) error {
	// Need session id
	return fiber.ErrNotImplemented
}

// Drops all user's refresh sessions (equivalent to logout)
func (s *Sessions) DropAll(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return ErrUserIDNotPassed
	}

	if err := s.r.DropAll(userID); err != nil {
		fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
