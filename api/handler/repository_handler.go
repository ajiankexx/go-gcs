package handler

import (
	"go-gcs/appError"
	"go-gcs/model"
	"go-gcs/service"
	"go-gcs/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

type RepoHandler struct {
	Service *service.RepoService
}

// @Summary      创建仓库
// @Description  当前登录用户创建一个新的代码仓库
// @Tags         Repository
// @Accept       json
// @Produce      json
// @Param        repoInfo  body  model.CreateRepoDTO  true  "仓库信息（名称、描述、是否私有、地址等）"
// @Param Authorization header string true "Authorization token (Bearer <token>)"
// @Router       /repository/create [post]
func (r *RepoHandler) CreateRepo(c *gin.Context) {
	req := &model.CreateRepoDTO{}
	err := utils.ValidateReq(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing id"})
		return
	}
	ctx := c.Request.Context()
	id, exists := userId.(int64)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "id type error"})
		return
	}
	err = r.Service.CreateRepo(ctx, req, id)
	if err != nil {
		if err == appError.ErrorRepoAlreadyExist {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "Create Repo Successful",
	})
	return
}

// @Summary      update repo
// @Description  update repo
// @Tags         Repository
// @Accept       json
// @Produce      json
// @Param        repoInfo  body  model.UpdateRepoDTO  true  "仓库信息（名称、描述、是否私有、地址等）"
// @Param Authorization header string true "Authorization token (Bearer <token>)"
// @Router       /repository/update [post]
func (r *RepoHandler) UpdateRepo(c *gin.Context) {
	req := &model.UpdateRepoDTO{}
	err := utils.ValidateReq(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing id"})
		return
	}

	ctx := c.Request.Context()
	id, exists := userId.(int64)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "id type error"})
		return
	}
	err = r.Service.UpdateRepo(ctx, req, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Update Repository Successful"})
}
