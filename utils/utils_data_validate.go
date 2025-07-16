package utils

import (
	"strings"
	"reflect"
	"context"
	"errors"

	"github.com/gin-gonic/gin"
)

type literal interface {
	int | int32 | int64 | float32 | float64 | string
}

func ValidateReq(c *gin.Context, tgt any) error {
	t := reflect.TypeOf(tgt)
	if t == nil || t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errors.New("TypeError")
	}
	return c.ShouldBindBodyWithJSON(tgt)
}

func ReadFromContext[T any](ctx context.Context, key any) (T, bool) {
	v := ctx.Value(key)
	val, ok := v.(T)
	return val, ok
}

func LiteralPtr[T literal](obj T) (*T){
	return &obj
}

func StructToMap(obj any, tag string) map[string]interface {} {
	out := make(map[string]interface{})
	if obj == nil {
		return out
	}

	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return out
		}
		v = v.Elem()
		t = t.Elem()
	}

	if v.Kind() != reflect.Struct {
		return out
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if fieldType.PkgPath != "" {
			continue
		}

		var key string
		if tag == "No" {
			key = fieldType.Name
		} else {
			tagValue := fieldType.Tag.Get(tag)
			if tagValue == "" {
				key = fieldType.Name
			} else {
				key = strings.Split(tagValue, ",")[0]
			}
		}

		var value interface{}
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				continue
			}
			value = field.Elem().Interface()
		} else {
			value = field.Interface()
		}

		out[key] = value
	}

	return out
}

