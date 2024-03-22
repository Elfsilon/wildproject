package controller

import model "temp/internal/app/domain/models"

type UsersService interface {
	Find(userID, email string) (model.User, error)
	FindDetailedByID(userID string) (model.UserDetailed, error)
	CountByEmail(email string) (int, error)
	Create(email, passwordHash string) (string, error)
}
