package service

import (
	"go-gcs/appError"
	"go-gcs/constants"
	"go-gcs/dao"
	"go-gcs/model"
	"go-gcs/mq"
	"go-gcs/utils"

	"time"
	"errors"
	"context"
	"encoding/json"
)

var (
	VerificationCodeTTL = constants.VerificationCodeTTL
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

//TODO: temporally only write how to send and process verify code request,
// however, currently not consider about what should be returned in this request.
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

func (r *UserService) SendVerificationCode(ctx context.Context, req *model.SendEmail) (error) {
	created_at := time.Now()
	expired_at := created_at.Add(VerificationCodeTTL)
	verificationCode := utils.GenVerifyCode()

	verifyCodeInfo := model.VerifiyCodeInfo{
		CreatedAt: created_at,
		ExpiredAt: expired_at,
		Code: verificationCode,
		Email: req.Email,
		Effective: true,
		Scene: "Unknown",
		// IP: req.IP,
	}

	data, err := json.Marshal(verifyCodeInfo)
	if err != nil {
		return err
	}
	redisKey := "verify_code:" + req.Email
	err = utils.GetRedisConn().Set(ctx, redisKey, data, VerificationCodeTTL).Err()
	if err != nil {
		return err
	}

	email_msg := &model.EmailMessage{
		Email: req.Email,
		Topic: "email-sender",
		Addr: "localhost:9092",
		VerificationCode: verificationCode,
	}
	email_sender := &mq.EmailSender{
		EmailMessage: email_msg,
	}
	err = email_sender.SendMessage()
	if err != nil {
		return err
	}
	return nil
}

func (r *UserService) UploadEmailAndVerifyCode(ctx context.Context, req *model.EmailAndVerifyCode) (error) {
	rdb := utils.GetRedisConn()
	redisKey := "verify_code:" + req.Email
	n, err := rdb.Exists(ctx, redisKey).Result()
	if err != nil {
		return err
	}
	if n == 0 {
		return appError.ErrorRedisNotFoundKey
	}
	val, err := rdb.Get(ctx, redisKey).Result()

	var verifyCodeInfo model.VerifiyCodeInfo
	err = json.Unmarshal([]byte(val), &verifyCodeInfo)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	if verifyCodeInfo.Code != req.VerifyCode {
		return appError.ErrorWrongVerifyCode
	}
	if time.Now().After(verifyCodeInfo.ExpiredAt) {
		return appError.ErrorExpiredVerifyCode
	}
	return nil
}
