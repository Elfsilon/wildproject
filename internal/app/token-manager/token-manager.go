package tokenmanager

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidSigningMethod   = errors.New("invalid token's signing method")
	ErrInvalidClaims          = errors.New("invalid claims")
	ErrClaimsEmptySessionID   = errors.New("claims' session_id is empty")
	ErrClaimsInvalidSessionID = errors.New("claims' session_id is invalid")
	ErrClaimsEmptyUserID      = errors.New("claims' user_id is empty")
)

type TokenManager struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenManager(secret []byte, ttl time.Duration) *TokenManager {
	return &TokenManager{secret, ttl}
}

func (tm *TokenManager) Generate(sessionID int, userID string) (string, error) {
	payload := jwt.RegisteredClaims{
		ID:        fmt.Sprint(sessionID),
		Subject:   userID,
		Issuer:    "wildproject",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.ttl).UTC()),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return jwtToken.SignedString(tm.secret)
}

func (tm *TokenManager) Validate(accessToken string) (int, string, error) {
	validator := func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidSigningMethod
		}

		return tm.secret, nil
	}

	var claims jwt.RegisteredClaims

	_, err := jwt.ParseWithClaims(accessToken, &claims, validator)
	if err != nil {
		return -1, "", err
	}

	rawSessionID := claims.ID
	if rawSessionID == "" {
		return -1, "", ErrClaimsEmptySessionID
	}

	sessionID, err := strconv.Atoi(rawSessionID)
	if err != nil {
		return -1, "", ErrClaimsInvalidSessionID
	}

	userID := claims.Subject
	if userID == "" {
		return -1, "", ErrClaimsEmptyUserID
	}

	return sessionID, userID, nil

}
