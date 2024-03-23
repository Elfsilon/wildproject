-- Active: 1710705105856@@127.0.0.1@5432@wildb-dev

DROP TABLE IF EXISTS user_info;
DROP TABLE IF EXISTS sex;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
  user_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email text NOT NULL UNIQUE,
  password_hash text NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
  updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp
);

CREATE TABLE sex (
  sex_id smallserial PRIMARY KEY,
  label varchar(16) NOT NULL UNIQUE
);

CREATE TABLE user_info (
  user_id uuid PRIMARY KEY 
    REFERENCES users (user_id)
      ON UPDATE CASCADE
      ON DELETE CASCADE,
  sex_id smallint DEFAULT 1 
    REFERENCES sex (sex_id)
      ON UPDATE CASCADE
      ON DELETE SET DEFAULT,
  name varchar(200) DEFAULT ''
);

-- CREATE TABLE user_role (
--   role_id smallserial PRIMARY KEY,
--   label varchar(16) NOT NULL UNIQUE
-- );

-- CREATE TABLE users_role (
--   user_id uuid PRIMARY KEY 
--     REFERENCES users (user_id)
--       ON UPDATE CASCADE
--       ON DELETE CASCADE,
--   sex_id smallint DEFAULT 1 
--     REFERENCES sex (sex_id)
--       ON UPDATE CASCADE
--       ON DELETE SET DEFAULT,
--   name varchar(200) DEFAULT ''
-- );

INSERT INTO sex (label) VALUES ('Не установлен'), ('Мужской'), ('Женский');

SELECT 
  u.user_id, 
  u.email, 
  u.password_hash, 
  u.created_at, 
  u.updated_at,
  ui.sex_id,
  ui.name 
FROM users AS u
INNER JOIN user_info AS ui
  USING(user_id)
WHERE u.user_id = 'c0459229-883b-4600-84ff-b0eecdabfbbd';

SELECT * FROM user_info;