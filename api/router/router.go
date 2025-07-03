package router

import (
	"go-gcs/constants"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

var ALL_API_PREFIX = constants.ALL_API_PREFIX

func RouterSetup() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	v1 := r.Group(ALL_API_PREFIX)
	{
		SetupUserRoutes(v1)
	}
	return r
}
