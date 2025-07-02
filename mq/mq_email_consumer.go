package mq

import (
	"go-gcs/model"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go-gcs/utils"
)

type EmailReader struct {
	EmailMessage *model.EmailMessage
}

func (r *EmailReader) ReadMessage() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic: "email-sender",
		Partition: 0,
		MaxBytes: 10e6,
	})

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			continue
		}

		var msg model.EmailMessage
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			continue
		}
		go utils.SendEmail(msg.Email, msg.VerificationCode)
	}
}
