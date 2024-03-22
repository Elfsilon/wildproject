package entity

type User struct {
	ID           string
	Email        string
	PasswordHash string
	CreatedAt    string
	UpdatetdAt   string
}

type UserDetailed struct {
	User
	SexID int
	Name  string
}
