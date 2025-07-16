package model

import (
	"time"
)

// temporarily don't add gmt_created, gmt_updated, gmt_deleted to UserDTO
type UserDO struct {
	userTableName
	Id        int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserName  string `gorm:"column:username" json:"username"`
	Email     string `gorm:"emeil" json:"email"`
	Password  string `gorm:"password" json:"password"`
	AvatarURL string `gorm:"avatar_url" json:"avarar_url"`
}

type UserDTO struct {
	userTableName
	UserName  string `gorm:"column:username" json:"username"`
	Email     string `gorm:"emeil" json:"email"`
	Password  string `gorm:"password" json:"password"`
	AvatarURL string `gorm:"avatar_url" json:"avarar_url"`
}

type UpdateUserDTO struct {
	userTableName
	UserName  *string `gorm:"column:username" json:"username"`
	Email     *string `gorm:"emeil" json:"email"`
	Password  *string `gorm:"password" json:"password"`
	AvatarURL *string `gorm:"avatar_url" json:"avarar_url"`
}
type UpdateUserDO  struct {
	UpdateUserDTO
}

type UserVO struct {
	userTableName
	UserName  string `gorm:"column:username" json:"username"`
	Email     string `gorm:"emeil" json:"email"`
	AvatarURL string `gorm:"avatar_url" json:"avarar_url"`
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
