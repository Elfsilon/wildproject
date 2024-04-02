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
			ui.name,
			ui.img_url
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

	UpdateUserName = `
		UPDATE user_info
			SET name = $1
		WHERE user_id = $2
		RETURNING name;
	`

	UpdateUserImgURL = `
		UPDATE user_info
			SET img_url = $1
		WHERE user_id = $2
		RETURNING img_url;
	`

	UpdateUserEmail = `
		UPDATE users
			SET email = $1
		WHERE user_id = $2
		RETURNING email;
	`

	UpdateUserSex = `
		UPDATE user_info
			SET sex_id = $1
		WHERE user_id = $2
		RETURNING sex_id;
	`

	UpdateUserPasswordHash = `
		UPDATE users
			SET password_hash = $1
		WHERE user_id = $2;
	`
)
