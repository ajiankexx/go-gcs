package model

import (
	"time"
)

// temporarily don't add gmt_created, gmt_updated, gmt_deleted to UserDTO
type UserDO struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AvatarURL string `json:"avarar_url"`
}

type UserDTO struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AvatarURL string `json:"avatar_url"`
}

type UpdateUserDTO struct {
	Username  *string `json:"username"`
	Email     *string `json:"email"`
	Password  *string `json:"password"`
	AvatarURL *string `json:"avatar_url"`
}

type UserVO struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type UpdatePasswordWithOldPasswordDTO struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type SendEmailDTO struct {
	Email string `json:"email"`
	// IP string `json:"ip"` //TODO: maybe changed later
}

type EmailAndVerifyCodeDTO struct {
	Email      string `json:"email"`
	VerifyCode string `json:"verify_code"`
}

type VerifiyCodeInfoDTO struct {
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
	Email     string    `json:"email"`
	Code      string    `json:"code"`
	Effective bool      `json:"effective"`
	// IP        string    `json:"ip,omitempty"`
	Scene string `json:"scene"`
}
