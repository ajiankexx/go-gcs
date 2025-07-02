package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"crypto/tls"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
	"github.com/jordan-wright/email"
)

func SendEmail(targetEmail, content string) error {
	err := godotenv.Load()
	if err != nil {
		log.Printf("加载.env文件失败: %v", err)
	}

	senderEmail := os.Getenv("SenderEmail")
	authCode := os.Getenv("EmailAuthorizationCode")
	smtpServer := os.Getenv("SenderServer")
	smtpPort := os.Getenv("SenderPort")

	e := email.NewEmail()
	e.From = senderEmail
	e.To = []string{targetEmail}
	e.Subject = "Authrization Code"
	e.Text = []byte(content)

	auth := smtp.PlainAuth(
		"",
		senderEmail,
		authCode,
		smtpServer,
	)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer,
	}

	err = e.SendWithTLS(smtpServer+":"+smtpPort, auth, tlsConfig)
	if err != nil {
		log.Printf("邮件发送失败: %v", err)
		return err
	}

	log.Println("✅ 邮件发送成功！")
	return nil
}

func GenVerifyCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
