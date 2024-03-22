package model

import (
	"temp/internal/stamp"
)

type (
	User struct {
		ID        string      `json:"user_id"`
		Email     string      `json:"email"`
		CreatedAt stamp.Stamp `json:"created_at,omitempty"`
		UpdatedAt stamp.Stamp `json:"updated_at,omitempty"`
	}

	UserDetailed struct {
		User
		Name string `json:"name,omitempty"`
		Sex  int    `json:"sex,omitempty"`
	}

	UserWithCredentials struct {
		User
		PasswordHash string `json:"-"`
	}
)
