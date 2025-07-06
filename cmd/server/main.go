package main

import (
	"go-gcs/api/router"
	"go-gcs/model"
	"go-gcs/mq"
	"go-gcs/logger"

	"go.uber.org/zap"

	_ "go-gcs/cmd/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title go-gcs
// @version 1.0
// @description This is go-gcs

// @host localhost:1234
// @BasePath /api/v1

// @schemes http // 支持的协议
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// email_message := &model.EmailMessage{Topic: "email-sender", Addr: "localhost:9200"}
	email_process := mq.EmailReader{EmailMessage: &model.EmailMessageDTO{}}
	go email_process.ReadMessage()

	logger.InitLogger()

	r := router.RouterSetup()
	r.GET("/swagger/*any", ginSwagger.WrapHandler((swaggerFiles.Handler)))
	port := ":1234"
	zap.L().Info("🚀 Server starting...", zap.String("port", port))
	r.Run(":1234")
}
