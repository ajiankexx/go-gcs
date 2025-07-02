package service

import (
	"context"
	"errors"
	"go-gcs/appError"
	"go-gcs/dao"
	"go-gcs/model"
	"go-gcs/utils"
	"go-gcs/mq"
)

type UserService struct {
	DAO *dao.UserDB
}

func (r *UserService) CreateUser(ctx context.Context, req *model.User) (*model.User, error) {
	req.Password = utils.Encrypt(req.Password)

	userData, err := r.DAO.GetUserByName(ctx, req.Username)
	if err == nil {
		return userData, appError.ErrorUserAlreadyExists
	} else if err != appError.ErrorUserNotFound {
		return userData, err
	}

	err = r.DAO.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *UserService) UpdateUser(ctx context.Context, req *model.UpdateUser, id string) (*model.User, error) {
	data, err := r.DAO.GetUserByID(ctx, id)
	original_data := data
	if err != nil {
		return nil, err
	}
	if req.Username != nil {
		if *req.Username == "" {
			return original_data, errors.New("username can't be empty")
		} else {
			data.Username = *req.Username
		}
	}

	if req.Email != nil {
		if *req.Email == "" {
			return original_data, errors.New("email can't be empty")
		} else {
			data.Email = *req.Email
		}
	}

	if req.AvatarURL != nil {
		if *req.AvatarURL == "" {
			return original_data, errors.New("AvatarURL can't be empty")
		} else {
			data.AvatarURL = *req.AvatarURL
		}
	}

	if req.Password != nil {
		if *req.Password == "" {
			return original_data, errors.New("Password can't be empty")
		} else {
		}
	}

	err = r.DAO.UpdateUser(ctx, data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *UserService) UpdatePasswordWithOldPassword(ctx context.Context, req *model.UpdatePasswordWithOldPassword, id string) (error) {
	user, err := r.DAO.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	if req.NewPassword == req.OldPassword {
		return appError.ErrorInvalidPassword
	}
	req.OldPassword = utils.Encrypt(req.OldPassword)
	if req.OldPassword != user.Password {
		return appError.ErrorWrongPassword
	}
	if req.NewPassword == "" {
		return appError.ErrorInvalidPassword
	}
	user.Password = utils.Encrypt(req.NewPassword)
	err = r.DAO.UpdateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserService) SendVerificationCode(ctx context.Context, email_msg *model.EmailMessage) (error) {

	email_sender := &mq.EmailSender{
		EmailMessage: email_msg,
	}
	err := email_sender.SendMessage()
	if err != nil {
		return err
	}
	return nil
}
