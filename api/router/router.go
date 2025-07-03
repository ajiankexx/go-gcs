package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func RouterSetup() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	v1 := r.Group("/api/v1")
	{
		SetupUserRoutes(v1)
	}
	return r
}
