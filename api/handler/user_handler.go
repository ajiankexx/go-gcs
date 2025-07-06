package handler

import (
	"go-gcs/model"
	"go-gcs/service"

	"net/http"
	"reflect"
	"errors"

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
// @Param user body model.UserDTO true "用户信息"
// @Success 200 {object} model.UserDTO "创建成功返回用户信息"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /users/create [post]
func (r *UserHandler) CreateUser(c *gin.Context) {
	req := &model.UserDTO{}
	err := r.ValidateUser(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	userVO, err := r.Service.CreateUser(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "user created",
		"user": userVO,
	})
}

// @Summary update user information
// @Description update user information
// @Tags Users
// @Router /users/update [post]
// @Param UpdateUser body model.UpdateUserDTO true "struct for update user"
// @Param Authorization header string true "Authorization token (Bearer <token>)"
func (r *UserHandler) UpdateUser(c *gin.Context) {
	req := &model.UpdateUserDTO{}
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
	userVO, err := r.Service.UpdateUser(ctx, req, idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "user updated",
		"user": userVO,
	})
}

// @Summary update password
// @Description update password with old password
// @Tags Users
// @Router /users/update-password-with-old-password [post]
// @Param UpdatePassword body model.UpdatePasswordWithOldPasswordDTO true "the struct for update password"
// @Param Authorization header string true "Authorization token (Bearer <token>)"
func (r *UserHandler) UpdatePasswordWithOldPassword(c *gin.Context) {
	req := &model.UpdatePasswordWithOldPasswordDTO{}
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
// @Param UpdatePasswordWithOldPassword body model.UpdatePasswordWithOldPasswordDTO  true "struct for update password"
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
// @Param request body model.SendEmailDTO true "the struct for Sending Verification Code"
func (r *UserHandler) SendVerificationCode(c *gin.Context) {
	req := &model.SendEmailDTO{}
	err := r.ValidateUser(c, req)
	// req.IP = c.ClientIP()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctx := c.Request.Context()
	err = r.Service.SendVerificationCode(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Status": "Email send successful"})
}

// @Summary upload email and verification code
// @Description veirify code
// @Tags Users
// @Router /users/upload-email-and-verifycode [post]
// @Param request body model.EmailAndVerifyCodeDTO true "the struct for upload email and verifycode"
func (r *UserHandler) UploadEmailAndVerifyCode(c *gin.Context) {
	req := &model.EmailAndVerifyCodeDTO{}
	err := r.ValidateUser(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctx := c.Request.Context()
	err = r.Service.UploadEmailAndVerifyCode(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Status": "Verify Successful"})
}
