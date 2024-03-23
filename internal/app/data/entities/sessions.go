package entity

type RefreshSession struct {
	SessionID    int
	RefreshToken string
	AccessToken  string
	UserID       string
	Uagent       string
	Fprint       string
	ExpiresAt    string
	CreatedAt    string
}
