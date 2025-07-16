package handler
import (
	"go-gcs/service"
	"github.com/gin-gonic/gin"
	"go-gcs/model"
	"go-gcs/utils"

	"net/http"
	"context"
)
type SshHandler struct {
	Service *service.SshKeyService
}

func (r *SshHandler) UploadSshKey(c *gin.Context) {
	req := &model.SshKeyDTO{}
	err := utils.ValidateReq(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	userIdRaw, _ := c.Get("userId")
	userId, ok := userIdRaw.(int64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong user id"})
		return
	}
	ctx = context.WithValue(ctx, "userId", userId)
	err = r.Service.UploadSshKey(ctx, req, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "sshKey upload successful"})
}

func (r *SshHandler) UpdateSshKey(c *gin.Context) {
	req := &model.UpdateSshKeyDTO{}
	err := utils.ValidateReq(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	userIdRaw, _ := c.Get("userId")
	userId, _ := userIdRaw.(int64)
	ctx = context.WithValue(ctx, "userId", userId)
	err = r.Service.UpdateSshKey(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "update ssh key successful"})
}

func (r *SshHandler) DeleteSshKey(c *gin.Context) {
	req := &model.SshKeyDTO{}
	err := utils.ValidateReq(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	userIdRaw, _ := c.Get("userId")
	userId, _ := userIdRaw.(int64)
	ctx = context.WithValue(ctx, "userId", userId)
	err = r.Service.DeleteSshKey(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": "delete ssh key successful"})
}
