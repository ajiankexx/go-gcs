package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"go-gcs/appError"
	"go-gcs/auth"
	"go-gcs/dao"
	"go-gcs/model"
	"go-gcs/utils"
	"go.uber.org/zap"
)

type LoginService struct {
	Validator *validator.Validate
	DAO       *dao.UserDB
}

func (r *LoginService) LoginVerifyPassword(ctx context.Context, req model.LoginRequestDTO) error {
	if err := r.Validator.Struct(req); err != nil {
		utils.LogValidationErrors(err)
		return appError.ErrorWrongFormatRequestData
	}
	user, err := r.DAO.GetUserByUserName(ctx, nil, req.UserName)
	if err != nil {
		zap.L().Error("LoginVerifyPassword() failed", zap.Error(err))
		return err
		// return appError.ErrUserNotFound
	}
	if user.Password != utils.Encrypt(req.PassWord) {
		return appError.ErrorPasswordInvalid
	}
	return nil
}

func (r *LoginService) Login(ctx context.Context, req *model.LoginRequestDTO) (string, error) {
	if err := r.Validator.Struct(req); err != nil {
		utils.LogValidationErrors(err)
		return "", appError.ErrorWrongFormatRequestData
	}
	err := r.LoginVerifyPassword(ctx, *req)
	if err != nil {
		zap.L().Error("Login() failed", zap.Error(err))
		return "", err
	}
	var id int64
	id, err = r.DAO.GetUserIdByUserName(ctx, nil, req.UserName)
	if err != nil {
		zap.L().Error("Login() failed", zap.Error(err))
		return "", err
	}
	token, err := auth.GenerateToken(req.UserName, id)
	return token, err
}
