package service

import model "temp/internal/app/domain/models"

type UsersService interface {
	Find(userID, email string) (model.User, error)
	FindDetailedByID(userID string) (model.UserDetailed, error)
	IsRegistered(email string) (bool, error)
	Create(email, passwordHash string) (string, error)
	Authenticate(email, password string) (string, error)
}

type SessionsService interface {
	Find(sessionID int) (model.ClientRefreshSession, error)
	FindAll(userID, uagent, fprint string) ([]model.ClientRefreshSession, error)
	Create(userID, uagent, fprint string) (model.TokenPair, error)
	Refresh(token, userID, uagent, fprint string) (model.TokenPair, error)
	Validate(sessionID int, accessToken, uagent, fprint string) error
	DropAll(userID, uagent, fprint string) error
	Drop(sessionID int) error
}
