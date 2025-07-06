package mq

import (
	"context"
	"encoding/json"
	"go-gcs/model"
	"go-gcs/utils"
	"time"

	"github.com/segmentio/kafka-go"
)

type EmailReader struct {
	EmailMessage *model.EmailMessageDTO
}

func (r *EmailReader) ReadMessage() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic: "email-sender",
		GroupID: "email-consumer",
		Partition: 0,
		MaxBytes: 10e6,
		CommitInterval: time.Second,
	})

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			continue
		}

		var msg model.EmailMessageDTO
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			continue
		}
		go utils.SendEmail(msg.Email, msg.VerificationCode)
	}
}
