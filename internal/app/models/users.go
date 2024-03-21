package m

type (
	User struct {
		ID    string `json:"user_id"`
		Email string `json:"email"`
	}

	UserWithTimeFields struct {
		User
		TimeFields
	}

	UserInfo struct {
		Name string `json:"name,omitempty"`
		Sex  int    `json:"sex,omitempty"`
	}

	UserDetailed struct {
		UserWithTimeFields
		UserInfo
	}

	UserCredentials struct {
		Email        string
		PasswordHash string
	}
)
