package model

import ()

type EmailMessage struct {
	Name             string `json:"name"`
	Email            string `json:"to"`
	Subject          string `json:"subject"`
	VerificationCode string `json:"verification_code"`
	Topic            string `json:"topic"`
	Addr             string `json:"address"`
}
