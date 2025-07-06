package service

import (
	"go-gcs/appError"
	"go-gcs/constants"
	"go-gcs/dao"
	"go-gcs/model"
	"go-gcs/mq"
	"go-gcs/utils"

	"context"
	"encoding/json"
	"errors"
	"time"
)

var (
	VerificationCodeTTL = constants.VerificationCodeTTL
)

type UserService struct {
	DAO *dao.UserDB
}

func (r *UserService) CreateUser(ctx context.Context, req *model.UserDTO) (*model.UserVO, error) {
	req.Password = utils.Encrypt(req.Password)

	userDTO, err := r.DAO.GetUserByName(ctx, req.Username)
	if err == nil {
		return utils.ConvertUserDTOToUserVO(userDTO), appError.ErrorUserAlreadyExists
	} else if err != appError.ErrorUserNotFound {
		return nil, err
	}

	err = r.DAO.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return utils.ConvertUserDTOToUserVO(req), nil
}

//TODO: temporally only write how to send and process verify code request,
// however, currently not consider about what should be returned in this request.

// maybe, in go, for front-end received data, we should use a stuct whose component var are pointer.
func (r *UserService) UpdateUser(ctx context.Context, req *model.UpdateUserDTO, id string) (*model.UserVO, error) {
	data, err := r.DAO.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	//when we use *model.UpdateUser for this operation, all keywords in it is a pointer, so we can
	//quickly distinguish which keywords is empty and which is not set by user

	//Below we assume all keywords can't be empty
	if req.Username != nil {
		if *req.Username == "" {
			return nil, errors.New("username can't be empty")
		} else {
			data.Username = *req.Username
		}
	}

	if req.Email != nil {
		if *req.Email == "" {
			return nil, errors.New("email can't be empty")
		} else {
			data.Email = *req.Email
		}
	}

	if req.AvatarURL != nil {
		if *req.AvatarURL == "" {
			return nil, errors.New("AvatarURL can't be empty")
		} else {
			data.AvatarURL = *req.AvatarURL
		}
	}

	if req.Password != nil {
		if *req.Password == "" {
			return nil, errors.New("Password can't be empty")
		} else {
		}
	}

	err = r.DAO.UpdateUser(ctx, data)
	if err != nil {
		return nil, err
	}
	return utils.ConvertUserDTOToUserVO(data), nil
}

// all data received from front-end maybe processed with a struct whose var is pointer, not value
func (r *UserService) UpdatePasswordWithOldPassword(ctx context.Context, req *model.UpdatePasswordWithOldPasswordDTO, id string) error {
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

func (r *UserService) SendVerificationCode(ctx context.Context, req *model.SendEmailDTO) error {
	created_at := time.Now()
	expired_at := created_at.Add(VerificationCodeTTL)
	verificationCode := utils.GenVerifyCode()

	verifyCodeInfo := model.VerifiyCodeInfoDTO{
		CreatedAt: created_at,
		ExpiredAt: expired_at,
		Code:      verificationCode,
		Email:     req.Email,
		Effective: true,
		Scene:     "Unknown",
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

	email_msg := &model.EmailMessageDTO{
		Email:            req.Email,
		Topic:            "email-sender",
		Addr:             "localhost:9092",
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

func (r *UserService) UploadEmailAndVerifyCode(ctx context.Context, req *model.EmailAndVerifyCodeDTO) error {
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

	var verifyCodeInfo model.VerifiyCodeInfoDTO
	err = json.Unmarshal([]byte(val), &verifyCodeInfo)
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
