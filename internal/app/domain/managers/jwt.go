package manager

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	model "wildproject/internal/app/domain/models"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidSigningMethod   = errors.New("invalid token's signing method")
	ErrInvalidClaims          = errors.New("invalid claims")
	ErrClaimsEmptySessionID   = errors.New("claims' session_id is empty")
	ErrClaimsInvalidSessionID = errors.New("claims' session_id is invalid")
	ErrClaimsEmptyUserID      = errors.New("claims' user_id is empty")
)

type JwtAccessManager struct {
	secret []byte
	ttl    time.Duration
}

func NewJwtManager(
	secret []byte,
	ttl time.Duration,
) *JwtAccessManager {
	return &JwtAccessManager{secret, ttl}
}

func (tm *JwtAccessManager) Generate(sessionID int, userID string) (string, error) {
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

func (tm *JwtAccessManager) getValidateFn() func(t *jwt.Token) (interface{}, error) {
	return func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidSigningMethod
		}

		return tm.secret, nil
	}
}

func (tm *JwtAccessManager) ParseAndValidate(accessToken string) (model.TokenPayload, error) {
	var claims jwt.RegisteredClaims

	_, err := jwt.ParseWithClaims(accessToken, &claims, tm.getValidateFn())
	if err != nil {
		return model.TokenPayload{}, err
	}

	return tm.parseClaims(claims)
}

// Retrieves token claims ignoring validation except of wrong signing method
func (tm *JwtAccessManager) Parse(accessToken string) (model.TokenPayload, error) {
	var claims jwt.RegisteredClaims

	_, err := jwt.ParseWithClaims(accessToken, &claims, tm.getValidateFn())
	if err != nil && errors.Is(err, ErrInvalidSigningMethod) {
		return model.TokenPayload{}, err
	}

	return tm.parseClaims(claims)
}

func (tm *JwtAccessManager) parseClaims(claims jwt.RegisteredClaims) (model.TokenPayload, error) {
	rawSessionID := claims.ID
	if rawSessionID == "" {
		return model.TokenPayload{}, ErrClaimsEmptySessionID
	}

	sessionID, err := strconv.Atoi(rawSessionID)
	if err != nil {
		return model.TokenPayload{}, ErrClaimsInvalidSessionID
	}

	userID := claims.Subject
	if userID == "" {
		return model.TokenPayload{}, ErrClaimsEmptyUserID
	}

	data := model.TokenPayload{
		SessionID: sessionID,
		UserID:    userID,
	}

	return data, nil
}
