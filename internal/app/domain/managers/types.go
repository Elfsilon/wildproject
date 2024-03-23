package manager

import model "temp/internal/app/domain/models"

type TokenManager interface {
	Parse(accessToken string) (model.TokenPayload, error)
	Generate(sessionID int, userID string) (string, error)
	ParseAndValidate(accessToken string) (model.TokenPayload, error)
}
