package m

type User struct {
	ID    string `json:"user_id"`
	Email string `json:"email"`
	TimeFields
}

type UserInfo struct {
	Name string `json:"name,omitempty"`
	Sex  int    `json:"sex,omitempty"`
}

type UserDetailed struct {
	User
	UserInfo
}
