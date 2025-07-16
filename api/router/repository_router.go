package router

import (
	"go-gcs/api/handler"
	"go-gcs/constants"
	"go-gcs/dao"
	"go-gcs/middleware"
	"go-gcs/service"
	"go-gcs/utils"

	"github.com/gin-gonic/gin"
)

var REPO_API_PROFIX = constants.REPO_API_PREFIX

func SetupRepositoryRoutes(r *gin.RouterGroup) {
	repoDAO := &dao.RepoDB{DB: utils.GetGormDB()}
	userDAO := &dao.UserDB{DB: utils.GetGormDB()}
	repoService := &service.RepoService{RepoDAO: repoDAO, UserDAO: userDAO}
	repoHandler := &handler.RepoHandler{Service: repoService}

	repo := r.Group(REPO_API_PROFIX)

	repo.Use(middleware.AuthMiddleware())
	repo.POST("create", repoHandler.CreateRepo)
	repo.POST("update", repoHandler.UpdateRepo)
}
