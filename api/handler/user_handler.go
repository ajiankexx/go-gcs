package handler

import (
	"errors"
	"go-gcs/model"
	"go-gcs/service"
	"go-gcs/utils"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *service.UserService
}

func (r *UserHandler) ValidateUser(c *gin.Context, tgt any) error {
	t := reflect.TypeOf(tgt)
	if t == nil || t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errors.New("TypeError")
	}
	return c.ShouldBindBodyWithJSON(tgt)
}


// Create godoc
// @Summary Create user
// @Description Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body model.User true "用户信息"
// @Success 200 {object} model.User "创建成功返回用户信息"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /users/create [post]
func (r *UserHandler) CreateUser(c *gin.Context) {
	req := &model.User{}
	err := r.ValidateUser(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	user, err := r.Service.CreateUser(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "user created",
		"user": user,
	})
}

// @Summary update user information
// @Description update user information
// @Tags Users
// @Router /users/update [post]
// @Param UpdateUser body model.UpdateUser true "struct for update user"
// @Param Authorization header string true "Authorization token (Bearer <token>)"
func (r *UserHandler) UpdateUser(c *gin.Context) {
	req := &model.UpdateUser{}
	err := r.ValidateUser(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	id, Ok := c.Get("id")
	if !Ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing id"})
		return
	}
	ctx :=  c.Request.Context()
	idStr, exists := id.(string)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "id type error"})
		return
	}
	user, err := r.Service.UpdateUser(ctx, req, idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "user updated",
		"user": user,
	})
}

// @Summary update password
// @Description update password with old password
// @Tags Users
// @Router /users/update-password-with-old-password [post]
// @Param Authorization header string true "Authorization token (Bearer <token>)"
func (r *UserHandler) UpdatePasswordWithOldPassword(c *gin.Context) {
	req := &model.UpdatePasswordWithOldPassword{}
	err := r.ValidateUser(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	id, Ok := c.Get("id")
	if !Ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing id"})
		return
	}
	ctx := c.Request.Context()
	idStr, exists := id.(string)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "id type error"})
		return
	}
	err = r.Service.UpdatePasswordWithOldPassword(ctx, req, idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Status":"Successful"})
}

// @Summary Reset password
// @Description update password with email verification code
// @Tags Users
// @Param UpdatePasswordWithOldPassword body model.UpdatePasswordWithOldPassword  true "struct for update password"
// @Router /users/update-password-with-email-verification-code [post]
func (r *UserHandler) UpdatePasswordWithEmailVerificationCode(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "Reset successful",
	})
}

// @Summary get an email verification code
// @Description get an email verification code
// @Tags Users
// @Router /users/get-email-verification-code [post]
// @Param request body model.SendEmail true "the struct for Sending Verification Code"
func (r *UserHandler) SendVerfificationCode(c *gin.Context) {
	req := &model.SendEmail{}
	err := r.ValidateUser(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctx := c.Request.Context()

	verificationCode, _ := utils.GenVerifyCode()
	email_msg := &model.EmailMessage{
		Email:req.Email, 
		Topic: "email-sender", 
		Addr: "localhost:9092",
		VerificationCode: verificationCode,
	}

	err = r.Service.SendVerificationCode(ctx, email_msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_here": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Status": "Email send successful"})
}
