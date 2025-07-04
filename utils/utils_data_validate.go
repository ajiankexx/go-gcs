package utils

import (
	"reflect"
	"errors"

	"github.com/gin-gonic/gin"
)

func ValidateReq(c *gin.Context, tgt any) error {
	t := reflect.TypeOf(tgt)
	if t == nil || t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errors.New("TypeError")
	}
	return c.ShouldBindBodyWithJSON(tgt)
}
