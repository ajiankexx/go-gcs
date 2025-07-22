
package router

import (
	"go-gcs/api/handler"
	"go-gcs/service"
	"go-gcs/dao"
	"go-gcs/middleware"
	"go-gcs/utils"
	"go-gcs/constants"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var SSH_KEY_API_PROFIX = constants.SSH_KEY_API_PROFIX

func SetupSshKeyRoutes(r *gin.RouterGroup) {
	sshKeyDAO := &dao.SshKeyDB{DB: utils.GetGormDB()}
	sshKeyService := &service.SshKeyService{SshKeyDAO: sshKeyDAO, Validator: validator.New()}
	sshHandler := &handler.SshHandler{Service: sshKeyService}

	ssh := r.Group(SSH_KEY_API_PROFIX)
	ssh.Use(middleware.AuthMiddleware())
	ssh.POST("upload", sshHandler.UploadSshKey)
	ssh.POST("update", sshHandler.UpdateSshKey)
	ssh.DELETE("delete", sshHandler.DeleteSshKey)
}
