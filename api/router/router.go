package router

import (
	"github.com/gin-gonic/gin"
)

func RouterSetup() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		SetupUserRoutes(v1)
	}
	return r
}
