package service

import (
	"database/sql"
	"errors"
	entity "temp/internal/app/data/entities"
	model "temp/internal/app/domain/models"
	"temp/internal/stamp"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrIDAndEmailEmpty = errors.New("expected user_id or email, but both are empty")
)

type Users struct {
	r UsersRepo
}

func NewUsers(r UsersRepo) *Users {
	return &Users{r}
}

// Find user either by passed id or by passed email
func (u *Users) Find(userID, email string) (model.User, error) {
	var ent entity.User
	var err error

	if userID != "" {
		ent, err = u.r.FindByID(userID)
	} else if email != "" {
		ent, err = u.r.FindByEmail(email)
	} else {
		return model.User{}, ErrIDAndEmailEmpty
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}

		return model.User{}, err
	}

	user := model.User{
		ID:        ent.ID,
		Email:     ent.Email,
		CreatedAt: stamp.Parse(ent.CreatedAt),
		UpdatedAt: stamp.Parse(ent.UpdatetdAt),
	}

	return user, nil
}

func (u *Users) FindDetailedByID(userID string) (model.UserDetailed, error) {
	ent, err := u.r.FindDetailedByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}

		return model.UserDetailed{}, err
	}

	user := model.UserDetailed{
		User: model.User{
			ID:        ent.ID,
			Email:     ent.Email,
			CreatedAt: stamp.Parse(ent.CreatedAt),
			UpdatedAt: stamp.Parse(ent.UpdatetdAt),
		},
		Name: ent.Name,
		Sex:  ent.SexID,
	}

	return user, nil
}

// Count
func (u *Users) IsRegistered(email string) (bool, error) {
	count, err := u.r.CountByEmail(email)
	if err != nil {
		return true, err
	}

	return count > 0, nil
}

// ...
func (u *Users) Create(email, password string) (string, error) {
	registered, err := u.IsRegistered(email)
	if err != nil {
		return "", err
	}

	if registered {
		return "", ErrAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return u.r.Create(email, string(passwordHash))
}
