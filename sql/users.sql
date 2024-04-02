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
  name varchar(200) NOT NULL DEFAULT '',
  img_url text NOT NULL DEFAULT ''
);

INSERT INTO sex (label) VALUES ('Не установлен'), ('Мужской'), ('Женский');
