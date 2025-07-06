package mq

import (
	"go-gcs/model"
	"go-gcs/appError"

	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type EmailSender struct {
	EmailMessage *model.EmailMessageDTO
}

func (r *EmailSender) SendMessage() error {
	msgBytes, err := json.Marshal(r.EmailMessage)
	if err != nil {
		return appError.ErrorEmailSend
		// log.Fatal("failed to marshal email message: ", err)
	}

	writer := &kafka.Writer{
		Addr:  kafka.TCP(r.EmailMessage.Addr),
		Topic: r.EmailMessage.Topic,
	}

	msg := kafka.Message{
		Key: []byte(r.EmailMessage.Email),
		Value: msgBytes,
	}
	err = writer.WriteMessages(context.Background(), msg)
	if err != nil {
		// return fmt.Errorf("failed to write message: %w", err)
		return err
		// return appError.ErrorEmailSend
	}
	return nil
}
