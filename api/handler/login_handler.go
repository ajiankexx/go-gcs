package handler

import (
	"go-gcs/model"
	"go-gcs/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	LoginService *service.LoginService
}

// Login godoc
// @Summary      用户登录
// @Description  用户登录接口，传入用户名和密码，返回 JWT Token
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Param        request body model.LoginRequestDTO true "登录请求体"
// @Success      200 {object} map[string]string "成功返回 token"
// @Failure      400 {object} map[string]string "请求参数格式错误"
// @Failure      401 {object} map[string]string "用户名或密码错误"
// @Router       /login [post]
func (r *LoginHandler) Login(c *gin.Context) {
	var req model.LoginRequestDTO
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	token, err := r.LoginService.Login(ctx, &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token":"Bearer " + token})
	return
}
