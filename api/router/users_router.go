package router

import (
	"go-gcs/api/handler"
	"go-gcs/middleware"
	"go-gcs/dao"
	"go-gcs/service"
	"go-gcs/utils"
	"go-gcs/constants"
	
	"github.com/gin-gonic/gin"
)

var USER_API_PREFIX = constants.USER_API_PREFIX
var LOGIN_API_PREFIX = constants.LOGIN_API_PREFIX

func SetupUserRoutes(r *gin.RouterGroup) {
	userDAO := &dao.UserDB{DB: utils.GetDBConn()}
	userService := &service.UserService{DAO: userDAO}
	userHandler := &handler.UserHandler{Service: userService}

	users := r.Group(USER_API_PREFIX)
	users.POST("create", userHandler.CreateUser)
	users.POST("get-email-verification-code", userHandler.SendVerificationCode)
	users.POST("upload-email-and-verifycode", userHandler.UploadEmailAndVerifyCode)

	authUsers := r.Group(USER_API_PREFIX)
	authUsers.Use(middleware.AuthMiddleware())
	{
		authUsers.POST("update", userHandler.UpdateUser)
		authUsers.POST("delete", handler.NotImplemented)
		authUsers.GET("get", handler.NotImplemented)
	}


	loginService := &service.LoginService{DAO: userDAO}
	loginHandler := &handler.LoginHandler{LoginService: loginService}

	login := r.Group(LOGIN_API_PREFIX)
	{
		login.POST("", loginHandler.Login)
	}

	register := r.Group("/register")
	{
		register.POST("", handler.NotImplemented)
	}
}
