package tokenmanager

import (
	"errors"
	"fmt"
	"strconv"
	model "temp/internal/app/domain/models"
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

func (tm *TokenManager) getValidateFn() func(t *jwt.Token) (interface{}, error) {
	return func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidSigningMethod
		}

		return tm.secret, nil
	}
}

func (tm *TokenManager) Validate(accessToken string) (model.TokenData, error) {
	var claims jwt.RegisteredClaims

	_, err := jwt.ParseWithClaims(accessToken, &claims, tm.getValidateFn())
	if err != nil {
		return model.TokenData{}, err
	}

	return tm.ParseClaims(claims)
}

// Retrieves token claims ignoring validation except of wrong signing method
func (tm *TokenManager) GetClaims(accessToken string) (model.TokenData, error) {
	var claims jwt.RegisteredClaims

	_, err := jwt.ParseWithClaims(accessToken, &claims, tm.getValidateFn())
	if err != nil && errors.Is(err, ErrInvalidSigningMethod) {
		return model.TokenData{}, err
	}

	return tm.ParseClaims(claims)
}

func (tm *TokenManager) ParseClaims(claims jwt.RegisteredClaims) (model.TokenData, error) {
	rawSessionID := claims.ID
	if rawSessionID == "" {
		return model.TokenData{}, ErrClaimsEmptySessionID
	}

	sessionID, err := strconv.Atoi(rawSessionID)
	if err != nil {
		return model.TokenData{}, ErrClaimsInvalidSessionID
	}

	userID := claims.Subject
	if userID == "" {
		return model.TokenData{}, ErrClaimsEmptyUserID
	}

	data := model.TokenData{
		SessionID: sessionID,
		UserID:    userID,
	}

	return data, nil
}
