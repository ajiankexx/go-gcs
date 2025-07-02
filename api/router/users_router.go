package router

import (
	"go-gcs/api/handler"
	"go-gcs/middleware"
	"go-gcs/dao"
	"go-gcs/service"
	"go-gcs/utils"
	
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.RouterGroup) {
	utils.InitDB("postgres://admin:1234@localhost:5432/gcs_db")

	userDAO := &dao.UserDB{DB: utils.GetDBConn()}
	userService := &service.UserService{DAO: userDAO}
	userHandler := &handler.UserHandler{Service: userService}

	users := r.Group("/users")
	users.POST("create", userHandler.CreateUser)
	users.POST("get-email-verification-code", userHandler.SendVerfificationCode)

	authUsers := r.Group("/users")
	authUsers.Use(middleware.AuthMiddleware())
	{
		authUsers.POST("update", userHandler.UpdateUser)
		authUsers.POST("delete", handler.NotImplemented)
		authUsers.GET("get", handler.NotImplemented)
	}


	loginService := &service.LoginService{DAO: userDAO}
	loginHandler := &handler.LoginHandler{LoginService: loginService}

	login := r.Group("/login")
	{
		login.POST("", loginHandler.Login)
	}

	register := r.Group("/register")
	{
		register.POST("", handler.NotImplemented)
	}
}
