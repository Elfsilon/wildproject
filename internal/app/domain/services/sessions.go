package service

import (
	"database/sql"
	"errors"
	entity "temp/internal/app/data/entities"
	repo "temp/internal/app/data/repositories"
	manager "temp/internal/app/domain/managers"
	model "temp/internal/app/domain/models"
	"temp/internal/stamp"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

var (
	ErrUnknownToken  = errors.New("unknown refresh token was used, your session dropped")
	ErrUnknownDevice = errors.New("unknown device was used, your session dropped")
	ErrExpiredToken  = errors.New("refresh token is expired, your session dropped")
)

type Sessions struct {
	cfg  *model.AuthConfig
	repo repo.SessionsRepo
	tm   manager.TokenManager
}

func NewSessions(
	cfg *model.AuthConfig,
	sr repo.SessionsRepo,
	tm manager.TokenManager,
) *Sessions {
	return &Sessions{cfg, sr, tm}
}

func (s *Sessions) Find(sessionID int) (model.ClientRefreshSession, error) {
	ent, err := s.repo.FindBySessionID(sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ClientRefreshSession{}, ErrNotFound
		}

		return model.ClientRefreshSession{}, err
	}

	session := model.ClientRefreshSession{
		SessionID: ent.SessionID,
		Uagent:    ent.Uagent,
		ExpiresAt: stamp.Parse(ent.ExpiresAt),
		CreatedAt: stamp.Parse(ent.CreatedAt),
	}

	return session, nil
}

// Find all sessions by "user_id" or find all by "user_id" and "device" if passed
//
// TODO: upgrade to pass { userID: required, uagent, fprint: optional }
func (s *Sessions) FindAll(userID, uagent, fprint string) ([]model.ClientRefreshSession, error) {
	var ents []entity.RefreshSession
	var err error

	if uagent != "" && fprint != "" {
		ents, err = s.repo.FindAllByDevice(userID, uagent, fprint)
	} else {
		ents, err = s.repo.FindAllByUserID(userID)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []model.ClientRefreshSession{}, ErrNotFound
		}

		return []model.ClientRefreshSession{}, err
	}

	sessions := make([]model.ClientRefreshSession, 0)
	for _, e := range ents {
		sessions = append(sessions, model.ClientRefreshSession{
			SessionID: e.SessionID,
			Uagent:    e.Uagent,
			ExpiresAt: stamp.Parse(e.ExpiresAt),
			CreatedAt: stamp.Parse(e.CreatedAt),
		})
	}

	return sessions, nil
}

// Drops all old user sessions associated with the device and creates new one
func (s *Sessions) Create(userID, uagent, fprint string) (model.TokenPair, error) {
	// TODO: figure out how to handle error
	s.DropAll(userID, uagent, fprint)

	return s.generateTokens(userID, uagent, fprint)
}

func (s *Sessions) Refresh(token, userID, uagent, fprint string) (model.TokenPair, error) {
	session, err := s.repo.FindByRefreshToken(token)
	if err != nil {
		return model.TokenPair{}, err
	}

	if err = s.repo.Drop(session.SessionID); err != nil {
		// TODO: add sentry logs
		log.Errorf("cannot drop session while refreshing token: %s", err)
	}

	if token != session.RefreshToken {
		return model.TokenPair{}, ErrUnknownToken
	}

	expiresAt := stamp.Parse(session.ExpiresAt).UTC()
	now := time.Now().UTC()

	log.Infof("Expires at: %v", expiresAt)
	log.Infof("Expires at: %v", expiresAt.Format(time.RFC1123Z))
	log.Infof("now: %v", now)
	log.Infof("now: %v", now.Format(time.RFC1123Z))
	log.Infof("expired? = %v", now.After(expiresAt))

	if now.After(expiresAt) {
		return model.TokenPair{}, ErrExpiredToken
	}

	return s.generateTokens(userID, uagent, fprint)
}

func (s *Sessions) generateTokens(userID, uagent, fprint string) (model.TokenPair, error) {
	rTokenExriresAt := time.Now().Add(s.cfg.RefreshTokenTTL).UTC()

	sessionID, refreshToken, err := s.repo.Create(userID, uagent, fprint, rTokenExriresAt)
	if err != nil {
		return model.TokenPair{}, err
	}

	accessToken, err := s.tm.Generate(sessionID, userID)
	if err != nil {
		return model.TokenPair{}, err
	}

	err = s.repo.SetAccessToken(sessionID, accessToken)
	if err != nil {
		return model.TokenPair{}, err
	}

	pair := model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return pair, nil
}

func (s *Sessions) Validate(sessionID int, accessToken, uagent, fprint string) error {
	old, err := s.repo.FindBySessionID(sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}

		return err
	}

	if accessToken != old.AccessToken {
		s.repo.Drop(old.SessionID)
		return ErrUnknownToken
	}

	if uagent != old.Uagent || fprint != old.Fprint {
		s.repo.Drop(old.SessionID)
		return ErrUnknownDevice
	}

	return nil
}

// Drops all user sessions associated with the device
func (s *Sessions) Drop(sessionID int) error {
	return s.repo.Drop(sessionID)
}

// Drops all user sessions associated with the device.
// Both "uagent" and "fprint" can be ommited
func (s *Sessions) DropAll(userID, uagent, fprint string) error {
	deviceSessions, err := s.FindAll(userID, uagent, fprint)
	if err != nil {
		return err
	}

	if deviceSessions != nil && len(deviceSessions) > 0 {
		// Try to delete sessions with this device
		for _, session := range deviceSessions {
			if err := s.repo.Drop(session.SessionID); err != nil {
				// TODO: Figure out what to do
				continue
			}
		}
	} else {
		// TODO: push notification: New device login detected
	}

	return nil
}
