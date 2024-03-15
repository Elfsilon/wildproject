DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS user_roles;

CREATE TABLE users (
  user_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email text NOT NULL UNIQUE,
  password_hash text NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
  updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp
);

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