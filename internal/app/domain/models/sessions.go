package model

import "wildproject/internal/stamp"

type RefreshSession struct {
	SessionID    int         `json:"session_id"`
	UserID       string      `json:"user_id,omitempty"`
	Uagent       string      `json:"user_agent"`
	Fprint       string      `json:"fingerprint"`
	RefreshToken string      `json:"refresh_token,omitempty"`
	AccessToken  string      `json:"access_token,omitempty"`
	ExpiresAt    stamp.Stamp `json:"expires_at,omitempty"`
}

type ClientRefreshSession struct {
	SessionID int         `json:"session_id"`
	Uagent    string      `json:"user_agent"`
	ExpiresAt stamp.Stamp `json:"expires_at,omitempty"`
	CreatedAt stamp.Stamp `json:"created_at,omitempty"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenPayload struct {
	SessionID int    `json:"session_id"`
	UserID    string `json:"usuer_id"`
}

type DeviceInfo struct {
	Uagent string `json:"user_agent"`
	Fprint string `json:"fingerprint"`
}

type CommonRequestPayload struct {
	TokenPayload
	DeviceInfo
}
