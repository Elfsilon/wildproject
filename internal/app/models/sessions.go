package m

import "github.com/golang-jwt/jwt/v5"

type RefreshSession struct {
	SessionID    int       `json:"session_id"`
	Uagent       string    `json:"user_agent"`
	Fprint       string    `json:"fingerprint"`
	UserID       string    `json:"user_id,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    Timestamp `json:"expires_at,omitempty"`
	CreatedAt    Timestamp `json:"created_at,omitempty"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenClaims struct {
	jwt.RegisteredClaims
	SessionID int    `json:"session_id"`
	UserID    string `json:"user_id"`
}
