package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	repo "temp/internal/app/data/repositories"
	model "temp/internal/app/domain/models"
	tokenmanager "temp/internal/app/token-manager"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidSessionID     = fiber.NewError(fiber.StatusBadRequest, "invalid session_id")
	ErrInvalidSession       = fiber.NewError(fiber.StatusBadRequest, "invalid locals session")
	ErrInvalidDevice        = fiber.NewError(fiber.StatusBadRequest, "invalid locals device_info")
	ErrInvalidTokenData     = fiber.NewError(fiber.StatusBadRequest, "invalid token_data")
	ErrInvalidAccessToken   = fiber.NewError(fiber.StatusBadRequest, "invalid access_token")
	ErrWrongEmailOrPassword = fiber.NewError(fiber.StatusBadRequest, "wrong email or password")
	ErrMismatchCurrentUser  = fiber.NewError(fiber.StatusInternalServerError, "current_user mismatches type of User")
	ErrAlreadyAuthorized    = fiber.NewError(fiber.StatusConflict, "passed token is valid and so user is already authorized")
	ErrUserAgentNotPassed   = fiber.NewError(fiber.StatusBadRequest, "User-Agent header is required")
	ErrFingerprintNotPassed = fiber.NewError(fiber.StatusBadRequest, "X-Fingerprint header is required")
	ErrUnknownRefreshToken  = fiber.NewError(fiber.StatusNotFound, "unknown refresh token was used, your session dropped")
	ErrExpiredRefreshToken  = fiber.NewError(fiber.StatusNotFound, "refresh token is expired, your session dropped")
	ErrSessionNotFound      = fiber.NewError(fiber.StatusNotFound, "session not found")
)

type Sessions struct {
	cfg *model.AuthConfig
	r   *repo.Sessions
	ur  *repo.Users
	tm  *tokenmanager.TokenManager
}

func NewSessions(
	cfg *model.AuthConfig,
	r *repo.Sessions,
	ur *repo.Users,
	tm *tokenmanager.TokenManager,
) *Sessions {
	return &Sessions{cfg, r, ur, tm}
}

func (s *Sessions) Get(c *fiber.Ctx) error {
	sessionID, err := strconv.Atoi(c.Params("session_id"))
	if err != nil {
		return ErrInvalidSessionID
	}

	session, err := s.r.FindBySessionID(sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrSessionNotFound
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(model.ClientRefreshSession{
		SessionID: session.SessionID,
		Uagent:    session.Uagent,
		ExpiresAt: session.ExpiresAt,
		CreatedAt: session.CreatedAt,
	})
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

	res := make([]model.ClientRefreshSession, 0)
	for _, s := range sessions {
		res = append(res, model.ClientRefreshSession{
			SessionID: s.SessionID,
			Uagent:    s.Uagent,
			ExpiresAt: s.ExpiresAt,
			CreatedAt: s.CreatedAt,
		})
	}

	return c.JSON(struct {
		Sessions []model.ClientRefreshSession `json:"sessions"`
	}{
		Sessions: res,
	})
}

func (s *Sessions) generateTokens(c *fiber.Ctx, userID, uagent, fprint string) error {
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

	return c.JSON(model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
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

	deviceSessions, err := s.r.FindByDevice(userID, uagent, fprint)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if deviceSessions != nil && len(deviceSessions) > 0 {
		// Try to delete sessions with this device
		for _, sessionID := range deviceSessions {
			if err := s.r.Drop(sessionID); err != nil {
				// TODO: Figure out what to do
				continue
			}
		}
	} else {
		// TODO: push notification: New device login detected
	}

	return s.generateTokens(c, userID, uagent, fprint)
}

// Generates a new pair of resfresh + access tokens by valid refresh token
func (s *Sessions) Refresh(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return ErrUserIDNotPassed
	}

	var payload struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&payload); err != nil {
		message := fmt.Sprintf("invalid body: %s", err)
		return fiber.NewError(fiber.StatusBadRequest, message)
	}

	oldSession, ok := c.Locals("session").(model.RefreshSession)
	if !ok {
		return ErrInvalidSession
	}

	device, ok := c.Locals("device_info").(model.DeviceInfo)
	if !ok {
		return ErrInvalidDevice
	}

	s.r.Drop(oldSession.SessionID)

	if payload.RefreshToken != oldSession.RefreshToken {
		return ErrUnknownRefreshToken
	}

	expiresAt := oldSession.ExpiresAt.UTC()
	now := time.Now().UTC()

	log.Infof("now: %v", now.Format(time.RFC822Z))
	log.Infof("expires: %v", expiresAt.Format(time.RFC822Z))
	log.Infof("expired?: %v", now.After(expiresAt))

	if now.After(expiresAt) {
		return ErrExpiredRefreshToken
	}

	return s.generateTokens(c, userID, device.Uagent, device.Fprint)
}

// Drops specified user's refresh session
func (s *Sessions) Drop(c *fiber.Ctx) error {
	sessionID, err := strconv.Atoi(c.Params("session_id"))
	if err != nil {
		return ErrInvalidSessionID
	}

	if err := s.r.Drop(sessionID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Need session id
	return c.SendStatus(fiber.StatusOK)
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
