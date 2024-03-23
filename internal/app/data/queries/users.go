package query

const (
	FindUserByID = `
		SELECT user_id, email, password_hash, created_at, updated_at
		FROM users 
		WHERE user_id = $1;
	`

	FindUserByEmail = `
		SELECT user_id, email, password_hash, created_at, updated_at
		FROM users 
		WHERE email = $1;
	`

	FindDetailedUserByID = `
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
		WHERE u.user_id = $1;
	`

	CountUsersByEmail = `
		SELECT count(*) FROM users WHERE email = $1;
	`

	CreateUser = `
		INSERT INTO users (email, password_hash) 
		VALUES ($1, $2) 
		RETURNING user_id;
	`

	CreateUserInfo = `
		INSERT INTO user_info (user_id) VALUES ($1);
	`
)
