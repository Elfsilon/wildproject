package rep

import (
	"database/sql"
	"errors"
	m "temp/internal/app/models"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrUserExists = errors.New("user already exists")
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{db}
}

func (u *Users) GetDetailed(userID string) (m.UserDetailed, error) {
	user, err := u.findByID(userID)
	if err != nil {
		return m.UserDetailed{}, err
	}
	detailed := m.UserDetailed{UserWithTimeFields: user}

	query := `
		SELECT name, sex_id
		FROM user_info
		WHERE user_id = $1;
	`

	err = u.db.QueryRow(query, userID).Scan(
		&detailed.Name, &detailed.Sex,
	)
	if err != nil {
		return m.UserDetailed{}, err
	}

	return detailed, nil
}

func (u *Users) Create(email, passwordHash string) (string, error) {
	count, err := u.countByEmail(email)
	if err != nil {
		return "", err
	}
	if count > 0 {
		return "", ErrUserExists
	}

	userID, err := u.createUser(email, passwordHash)
	if err != nil {
		return "", err
	}

	if err := u.createUserInfo(userID); err != nil {
		return "", err
	}

	return userID, nil
}

func (u *Users) createUser(email, passwordHash string) (string, error) {
	query := `
		INSERT INTO users (email, password_hash) 
		VALUES ($1, $2) 
		RETURNING user_id;
	`
	var userID string

	err := u.db.QueryRow(query, email, passwordHash).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (u *Users) createUserInfo(userID string) error {
	query := `INSERT INTO user_info (user_id) VALUES ($1);`

	err := u.db.QueryRow(query, userID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (u *Users) Update(userID string, info *m.UserInfo) error {
	return fiber.ErrNotImplemented
}

func (u *Users) Delete(userID string) error {
	return fiber.ErrNotImplemented
}

func (u *Users) ChangeEmail(userID, email string) error {
	return fiber.ErrNotImplemented
}

func (u *Users) findByID(userID string) (m.UserWithTimeFields, error) {
	query := `
		SELECT user_id, email, created_at, updated_at
		FROM users 
		WHERE user_id = $1;
	`
	var user m.UserWithTimeFields

	err := u.db.QueryRow(query, userID).Scan(
		&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return m.UserWithTimeFields{}, err
	}

	return user, nil
}

func (u *Users) findByEmail(email string) (m.UserWithTimeFields, error) {
	query := `
		SELECT user_id, email, created_at, updated_at
		FROM users 
		WHERE email = $1;
	`
	var user m.UserWithTimeFields

	err := u.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return m.UserWithTimeFields{}, err
	}

	return user, nil
}

func (u *Users) FindCredentialsByEmail(email string) (string, m.UserCredentials, error) {
	query := `
		SELECT user_id, email, password_hash
		FROM users 
		WHERE email = $1;
	`
	var userID string
	var credentials m.UserCredentials

	err := u.db.QueryRow(query, email).Scan(
		&userID, &credentials.Email, &credentials.PasswordHash,
	)
	if err != nil {
		return "", m.UserCredentials{}, err
	}

	return userID, credentials, nil
}

func (u *Users) countByEmail(email string) (int, error) {
	query := `SELECT count(*) FROM users WHERE email = $1;`

	var count int

	err := u.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
