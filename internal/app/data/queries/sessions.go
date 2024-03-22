package query

const (
	FindSessionsByUserID = `
		SELECT session_id, 
			refresh_token, 
			access_token, 
			user_id, 
			user_agent, 
			fingerprint, 
			expires_at, 
			created_at
		FROM refresh_sessions 
		WHERE user_id = $1;
	`

	FindSessionsByUserDevice = `
		SELECT session_id, 
			refresh_token, 
			access_token, 
			user_id, 
			user_agent, 
			fingerprint, 
			expires_at, 
			created_at
		FROM refresh_sessions 
		WHERE user_id = $1 
			AND user_agent = $2 
			AND fingerprint = $3;
	`

	FindSessionByID = `
		SELECT session_id, 
			refresh_token, 
			access_token, 
			user_id, 
			user_agent, 
			fingerprint, 
			expires_at, 
			created_at
		FROM refresh_sessions 
		WHERE session_id = $1;
	`

	CreateSession = `
		INSERT INTO refresh_sessions (user_id, user_agent, fingerprint, expires_at) 
		VALUES ($1, $2, $3, $4)
		RETURNING session_id, refresh_token;
	`

	SetSessionAccessToken = `
		UPDATE refresh_sessions
		SET access_token = $1
		WHERE session_id = $2;
	`

	DropSession = `
		DELETE FROM refresh_sessions
		WHERE session_id = $1;
	`

	DropAllSessions = `
		DELETE FROM refresh_sessions
		WHERE user_id = $1;
	`
)
