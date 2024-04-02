package service

import (
	"database/sql"
	"errors"
	entity "wildproject/internal/app/data/entities"
	repo "wildproject/internal/app/data/repositories"
	model "wildproject/internal/app/domain/models"
	"wildproject/internal/stamp"

	"golang.org/x/crypto/bcrypt"
)

const (
	passwordCost = 12
)

var (
	ErrIDAndEmailEmpty   = errors.New("expected user_id or email, but both are empty")
	ErrPasswordsMismatch = errors.New("passwords mismatches")
)

type Users struct {
	repo repo.UsersRepo
}

func NewUsers(r repo.UsersRepo) *Users {
	return &Users{r}
}

// Find user either by passed id or by passed email
func (u *Users) Find(userID, email string) (model.User, error) {
	var ent entity.User
	var err error

	if userID != "" {
		ent, err = u.repo.FindByID(userID)
	} else if email != "" {
		ent, err = u.repo.FindByEmail(email)
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
	ent, err := u.repo.FindDetailedByID(userID)
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
		Name:     ent.Name,
		Sex:      ent.SexID,
		ImageURL: ent.ImageURL,
	}

	return user, nil
}

func (u *Users) IsRegistered(email string) (bool, error) {
	count, err := u.repo.CountByEmail(email)
	if err != nil {
		return true, err
	}

	return count > 0, nil
}

func (u *Users) Create(email, password string) (string, error) {
	registered, err := u.IsRegistered(email)
	if err != nil {
		return "", err
	}

	if registered {
		return "", ErrAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)
	if err != nil {
		return "", err
	}

	return u.repo.Create(email, string(passwordHash))
}

func (u *Users) Authenticate(email, password string) (string, error) {
	user, err := u.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}

		return "", err
	}

	bHash := []byte(user.PasswordHash)
	bPass := []byte(password)

	if err := bcrypt.CompareHashAndPassword(bHash, bPass); err != nil {
		return "", ErrPasswordsMismatch
	}

	return user.ID, nil
}

func (u *Users) ChangeName(userID, name string) (string, error) {
	err := u.repo.ChangeName(userID, name)
	if err != nil {
		return "", err
	}

	return name, err
}

func (u *Users) ChangeSex(userID string, sexID int) (int, error) {
	err := u.repo.ChangeSex(userID, sexID)
	if err != nil {
		return -1, err
	}

	return sexID, err
}

func (u *Users) ChangeEmail(userID, email string) (string, error) {
	registered, err := u.IsRegistered(email)
	if err != nil {
		return "", err
	}

	if registered {
		return "", ErrAlreadyExists
	}

	err = u.repo.ChangeEmail(userID, email)
	if err != nil {
		return "", err
	}

	return email, err
}

func (u *Users) ChangePassword(userID, password string) error {
	phash, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)
	if err != nil {
		return err
	}

	return u.repo.ChangePasswordHash(userID, string(phash))
}

func (u *Users) ChangeImageURL(userID, url string) error {
	return u.repo.ChangeImageURL(userID, url)
}
