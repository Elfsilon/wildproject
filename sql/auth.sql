DROP TABLE IF EXISTS refresh_sessions;

CREATE TABLE refresh_sessions (
  refresh_token uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL 
    REFERENCES users (user_id)
      ON DELETE CASCADE
      ON UPDATE CASCADE,
  user_agent text NOT NULL,
  fingerprint text NOT NULL,
  expires_at timestamp with time zone NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT current_timestamp
);