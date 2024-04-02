package repo

import (
	"database/sql"
	"errors"
	"time"
	"wildproject/internal/app/data/database"
	entity "wildproject/internal/app/data/entities"
	query "wildproject/internal/app/data/queries"
)

type Sessions struct {
	db database.Instance
}

func NewSessions(db database.Instance) *Sessions {
	return &Sessions{db}
}

func (s *Sessions) FindAllByUserID(userID string) ([]entity.RefreshSession, error) {
	rows, err := s.db.Query(query.FindSessionsByUserID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []entity.RefreshSession{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	sessions := make([]entity.RefreshSession, 0)

	for rows.Next() {
		var session entity.RefreshSession

		err := rows.Scan(
			&session.SessionID, &session.RefreshToken, &session.AccessToken,
			&session.UserID, &session.Uagent, &session.Fprint, &session.ExpiresAt,
			&session.CreatedAt,
		)

		if err != nil {
			continue
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (s *Sessions) FindAllByDevice(
	userID, uagent, fprint string,
) (
	[]entity.RefreshSession, error,
) {
	sessions := make([]entity.RefreshSession, 0)

	rows, err := s.db.Query(query.FindSessionsByUserDevice, userID, uagent, fprint)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var session entity.RefreshSession

		err := rows.Scan(
			&session.SessionID, &session.RefreshToken, &session.AccessToken,
			&session.UserID, &session.Uagent, &session.Fprint, &session.ExpiresAt,
			&session.CreatedAt,
		)

		if err != nil {
			continue
		}

		sessions = append(sessions, session)
	}

	if len(sessions) == 0 {
		return []entity.RefreshSession{}, sql.ErrNoRows
	}

	return sessions, nil
}

func (s *Sessions) FindBySessionID(sessionID int) (entity.RefreshSession, error) {
	var session entity.RefreshSession

	err := s.db.QueryRow(query.FindSessionByID, sessionID).Scan(
		&session.SessionID, &session.RefreshToken, &session.AccessToken,
		&session.UserID, &session.Uagent, &session.Fprint, &session.ExpiresAt,
		&session.CreatedAt,
	)

	if err != nil {
		return entity.RefreshSession{}, err
	}

	return session, nil
}

func (s *Sessions) FindByRefreshToken(token string) (entity.RefreshSession, error) {
	var session entity.RefreshSession

	err := s.db.QueryRow(query.FindSessionByRefreshToken, token).Scan(
		&session.SessionID, &session.RefreshToken, &session.AccessToken,
		&session.UserID, &session.Uagent, &session.Fprint, &session.ExpiresAt,
		&session.CreatedAt,
	)

	if err != nil {
		return entity.RefreshSession{}, err
	}

	return session, nil
}

// Returns "session_id" and "refresh_token"
func (s *Sessions) Create(
	userID, uagent, fprint string, expriresAt time.Time,
) (
	int, string, error,
) {
	var sessionID int
	var refreshToken string

	err := s.db.QueryRow(
		query.CreateSession, userID, uagent, fprint, expriresAt,
	).Scan(&sessionID, &refreshToken)

	if err != nil {
		return -1, "", err
	}

	return sessionID, refreshToken, nil
}

func (s *Sessions) SetAccessToken(sessionID int, accessToken string) error {
	return s.db.QueryRow(query.SetSessionAccessToken, accessToken, sessionID).Err()
}

func (s *Sessions) Drop(sessionID int) error {
	return s.db.QueryRow(query.DropSession, sessionID).Err()
}

func (s *Sessions) DropAll(userID string) error {
	return s.db.QueryRow(query.DropAllSessions, userID).Err()
}
