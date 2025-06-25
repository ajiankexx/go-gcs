package main

import (
	"log"
	"go-gcs/api/router"

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
	r := router.RouterSetup()
	r.GET("/swagger/*any", ginSwagger.WrapHandler((swaggerFiles.Handler)))
	port := ":1234"
	log.Printf("Server starting on port http://localhost%s/api/v1/users", port)
	r.Run(":1234")
}
