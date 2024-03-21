-- Active: 1710705105856@@127.0.0.1@5432@wildb-dev
DROP TABLE IF EXISTS refresh_sessions;

CREATE TABLE refresh_sessions (
  session_id serial PRIMARY KEY,
  refresh_token uuid DEFAULT gen_random_uuid(),
  access_token text NOT NULL DEFAULT '',
  user_id uuid NOT NULL 
    REFERENCES users (user_id)
      ON DELETE CASCADE
      ON UPDATE CASCADE,
  user_agent text NOT NULL,
  fingerprint text NOT NULL,
  expires_at timestamp with time zone NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT current_timestamp
);

SELECT * FROM users;
SELECT * FROM refresh_sessions;

INSERT INTO refresh_sessions(user_id, user_agent, fingerprint, expires_at) VALUES ();

DELETE FROM refresh_sessions WHERE user_id = 'bcbbf682-d9bd-4aeb-9921-f1f14bb7e3fa';

DELETE FROM refresh_sessions
		WHERE session_id = 3;

SELECT session_id 
FROM refresh_sessions 
WHERE user_id = 'e7ec7a39-f2b9-4719-96f7-c30df814aa46'
  AND user_agent = 'PostmanRuntime/7.37.0' 
  AND fingerprint = 'Test fingerprint';