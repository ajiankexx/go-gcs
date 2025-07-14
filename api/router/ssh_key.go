package router

import (
	"go-gcs/api/handler"
	"go-gcs/service"
	"go-gcs/dao"
	"go-gcs/middleware"
	"go-gcs/utils"
	"go-gcs/constants"

	"github.com/gin-gonic/gin"
)

var SSH_KEY_API_PROFIX = constants.SSH_KEY_API_PROFIX

func SetupSshKeyRoutes(r *gin.RouterGroup) {
	sshDAO := &dao.SshDB{DB: utils.GetDBPool()}
	sshService := &service.SshService{DAO: sshDAO}
	sshHandler := &handler.SshHandler{Service: sshService}

	ssh := r.Group(SSH_KEY_API_PROFIX)
	ssh.Use(middleware.AuthMiddleware())
	ssh.POST("upload", sshHandler.UploadSsh)
	ssh.POST("update", sshHandler.UpdateSsh)
	ssh.GET("ssh-key-publickey", sshHandler.GetSshPublicKey)
	ssh.GET("ssh-key-name", sshHandler.GetSshKeyName)
	ssh.GET("page", sshHandler.Page)
	ssh.DELETE("delete", sshHandler.DeleteSshKey)
}
