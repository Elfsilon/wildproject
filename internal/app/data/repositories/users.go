package repo

import (
	"errors"
	"wildproject/internal/app/data/database"
	entity "wildproject/internal/app/data/entities"
	query "wildproject/internal/app/data/queries"
)

var (
	ErrUserExists = errors.New("user already exists")
)

type Users struct {
	db database.Instance
}

func NewUsers(db database.Instance) *Users {
	return &Users{db}
}

func (u *Users) FindByID(userID string) (entity.User, error) {
	var user entity.User

	err := u.db.QueryRow(query.FindUserByID, userID).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatetdAt,
	)

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u *Users) FindByEmail(email string) (entity.User, error) {
	var user entity.User

	err := u.db.QueryRow(query.FindUserByEmail, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatetdAt,
	)

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u *Users) FindDetailedByID(userID string) (entity.UserDetailed, error) {
	var user entity.UserDetailed

	err := u.db.QueryRow(query.FindDetailedUserByID, userID).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt,
		&user.UpdatetdAt, &user.SexID, &user.Name,
	)

	if err != nil {
		return entity.UserDetailed{}, err
	}

	return user, nil
}

func (u *Users) CountByEmail(email string) (int, error) {
	var count int

	err := u.db.QueryRow(query.CountUsersByEmail, email).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Creates user with corresponding user_info
//
// Returns "user_id" of created user or "" in error case
func (u *Users) Create(email, passwordHash string) (string, error) {
	var userID string

	row := u.db.QueryRow(query.CreateUser, email, passwordHash)
	err := row.Scan(&userID)

	if err != nil {
		return "", err
	}

	err = u.db.QueryRow(query.CreateUserInfo, userID).Err()
	if err != nil {
		return "", err
	}

	return userID, nil
}
