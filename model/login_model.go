package model

type LoginRequestDTO struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}
