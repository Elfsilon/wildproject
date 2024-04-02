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
