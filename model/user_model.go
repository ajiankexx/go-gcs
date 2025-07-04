package model

import (
	"time"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AvatarURL string `json:"avarar_url"`
}

type UpdateUser struct {
	Username  *string `json:"username, omitempty"`
	Email     *string `json:"email, omitempty"`
	Password  *string `json:"password, omitempty"`
	AvatarURL *string `json:"avatar_url, omitempty"`
}

type UpdatePasswordWithOldPassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type SendEmail struct {
	Email string `json:"email"`
	// IP string `json:"ip"` //TODO: maybe changed later
}

type EmailAndVerifyCode struct {
	Email      string `json:"email"`
	VerifyCode string `json:"verify_code"`
}

type VerifiyCodeInfo struct {
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
	Email     string    `json:"email"`
	Code      string    `json:"code"`
	Effective bool      `json:"effective"`
	// IP        string    `json:"ip,omitempty"`
	Scene string `json:"scene"`
}
