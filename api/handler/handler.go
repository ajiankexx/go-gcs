package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json: "email"`
}

// NotImplemented returns a standard response for unimplemented routes
// @Summary Not implemented
// @Description This endpoint is not yet implemented
// @Tags Unimplemented
// @Produce json
// @Success 501 {object} map[string]string "Not implemented"
func NotImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, map[string]string{"error":"This endpoint is not implemented yet"})
}
