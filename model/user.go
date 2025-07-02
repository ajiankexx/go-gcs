package model

type User struct {
	ID int64 `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	AvatarURL string `json:"avarar_url"`
}

type UpdateUser struct {
	Username *string `json:"username, omitempty"`
	Email *string `json:"email, omitempty"`
	Password *string `json:"password, omitempty"`
	AvatarURL *string `json:"avatar_url, omitempty"`
}

type UpdatePasswordWithOldPassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type SendEmail struct {
	Email string `json:"email"`
}
