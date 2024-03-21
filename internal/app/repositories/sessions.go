package rep

import (
	"database/sql"
	"errors"
	m "temp/internal/app/models"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type Sessions struct {
	db *sql.DB
}

func NewSessions(db *sql.DB) *Sessions {
	return &Sessions{db}
}

func (s *Sessions) FindByDevice(userID, uagent, fprint string) ([]int, error) {
	query := `
		SELECT session_id 
		FROM refresh_sessions 
		WHERE user_id = $1 
			AND user_agent = $2 
			AND fingerprint = $3;
	`
	ids := make([]int, 0)

	rows, err := s.db.Query(query, userID, uagent, fprint)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			continue
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (s *Sessions) FindByUserID(userID string) ([]m.RefreshSession, error) {
	query := `
		SELECT session_id, user_agent, fingerprint
		FROM refresh_sessions 
		WHERE user_id = $1;
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []m.RefreshSession{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	sessions := make([]m.RefreshSession, 0)

	for rows.Next() {
		var session m.RefreshSession

		err := rows.Scan(&session.SessionID, &session.Uagent, &session.Fprint)
		if err != nil {
			continue
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (s *Sessions) FindBySessionID(sessionID int) (m.RefreshSession, error) {
	query := `
		SELECT 
			session_id, refresh_token, access_token, user_id, 
			user_agent, fingerprint, expires_at, created_at
		FROM refresh_sessions 
		WHERE session_id = $1;
	`
	var session m.RefreshSession

	err := s.db.QueryRow(query, sessionID).Scan(
		&session.SessionID, &session.RefreshToken, &session.AccessToken, &session.UserID,
		&session.Uagent, &session.Fprint, &session.ExpiresAt, &session.CreatedAt,
	)

	if err != nil {
		return m.RefreshSession{}, err
	}

	return session, nil
}

func (s *Sessions) Create(
	userID, uagent, fprint string,
	expriresAt time.Time,
) (int, string, error) {
	query := `
		INSERT INTO refresh_sessions (user_id, user_agent, fingerprint, expires_at) 
		VALUES ($1, $2, $3, $4)
		RETURNING session_id, refresh_token;
	`
	var sessionID int
	var refreshToken string

	err := s.db.QueryRow(
		query, userID, uagent, fprint, expriresAt,
	).Scan(&sessionID, &refreshToken)

	if err != nil {
		return -1, "", err
	}

	return sessionID, refreshToken, nil
}

func (s *Sessions) SetAccessToken(sessionID int, accessToken string) error {
	query := `
		UPDATE refresh_sessions
		SET access_token = $1
		WHERE session_id = $2;
	`

	return s.db.QueryRow(query, accessToken, sessionID).Err()
}

func (s *Sessions) Drop(sessionID int) error {
	log.Infof("Drop %v", sessionID)
	query := `
		DELETE FROM refresh_sessions
		WHERE session_id = $1;
	`
	return s.db.QueryRow(query, sessionID).Err()
}

func (s *Sessions) DropAll(userID string) error {
	query := `
		DELETE FROM refresh_sessions
		WHERE user_id = $1;
	`
	return s.db.QueryRow(query, userID).Err()
}
