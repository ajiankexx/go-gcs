package model

type LoginRequestDTO struct {
	UserName string `json:"username" validate:"required,min=1,max=50"`
	PassWord string `json:"password" validate:"required,min=6"`
}
