package model

import ()

type EmailMessageDTO struct {
	Name             string `json:"name"               validate:"required,min=1,max=50"`
	Email            string `json:"to"                 validate:"required,email"`
	Subject          string `json:"subject"            validate:"required,min=1,max=100"`
	VerificationCode string `json:"verification_code"  validate:"required,len=6,numeric"`
	Topic            string `json:"topic"              validate:"required,min=1,max=50"`
	Addr             string `json:"address"            validate:"omitempty,max=255"`
}
