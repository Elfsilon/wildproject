package service

import (
	entity "temp/internal/app/data/entities"
	"time"
)

type SessionsRepo interface {
	FindAllByUserID(userID string) ([]entity.RefreshSession, error)
	FindAllByDevice(userID, uagent, fprint string) ([]entity.RefreshSession, error)
	FindBySessionID(sessionID int) (entity.RefreshSession, error)
	Create(userID, uagent, fprint string, expriresAt time.Time) (int, string, error)
	SetAccessToken(sessionID int, accessToken string) error
	Drop(sessionID int) error
	DropAll(userID string) error
}

type UsersRepo interface {
	FindByID(userID string) (entity.User, error)
	FindByEmail(email string) (entity.User, error)
	FindDetailedByID(userID string) (entity.UserDetailed, error)
	CountByEmail(email string) (int, error)
	Create(email, passwordHash string) (string, error)
}
