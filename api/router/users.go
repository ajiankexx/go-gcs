package router

import (
	"go-gcs/api/handler"
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
	{
		users.POST("create", userHandler.CreateUser)
		users.POST("update", userHandler.UpdateUser)
		users.POST("delete", handler.NotImplemented)
		users.GET("get", handler.NotImplemented)
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
