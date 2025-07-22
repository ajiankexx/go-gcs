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
	UserName  string `json:"username" validate:"required,min=1,max=50"`
	Email     string `json:"email" validate:"required,email,max=50"`
	Password  string `json:"password" validate:"required,min=6"`
	AvatarURL string `json:"avatar_url" validate:"required,url,max=1024"`
}

type UpdateUserDTO struct {
	userTableName
	UserName  *string `json:"username" validate:"omitempty,min=1,max=50"`
	Email     *string `json:"email" validate:"omitempty,email,max=50"`
	Password  *string `json:"password" validate:"omitempty,min=6"`
	AvatarURL *string `json:"avarar_url" validate:"omitempty,url,max=1024"`
}

type UpdateUserDO struct {
	UpdateUserDTO
}

type UserVO struct {
	userTableName
	UserName  string `gorm:"column:username" json:"username"`
	Email     string `gorm:"emeil" json:"email"`
	AvatarURL string `gorm:"avatar_url" json:"avarar_url"`
}

type UpdatePasswordWithOldPasswordDTO struct {
	OldPassword string `json:"old_password" validate:"required,min=6"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type SendEmailDTO struct {
	Email string `json:"email", validate:"required,email,max=50"`
	// IP string `json:"ip"` //TODO: maybe changed later
}

type EmailAndVerifyCodeDTO struct {
	Email      string `json:"email" validate:"required,email,max=50"`
	VerifyCode string `json:"verify_code" validate:"required,len=6"`
}

type VerifiyCodeInfoDTO struct {
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
	Email     string    `json:"email" validate:"required,email,max=50"`
	Code      string    `json:"code" validate:"required,len=6"`
	Effective bool      `json:"effective" validate:"required"`
	// IP        string    `json:"ip,omitempty"`
	Scene string `json:"scene"`
}
