package service

import (
	"go-gcs/appError"
	"go-gcs/constants"
	"go-gcs/dao"
	"go-gcs/model"
	"go-gcs/mq"
	"go-gcs/utils"
	"github.com/jinzhu/copier"

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

	userDO, err := r.DAO.GetUserByUserName(ctx, req.UserName)
	if err == nil {
		var userVO model.UserVO
		copier.Copy(userVO, userDO)
		return &userVO, appError.ErrorUserAlreadyExists
	} else if err != appError.ErrorUserNotFound {
		return nil, err
	}
	userDO = &model.UserDO{}
	copier.Copy(userDO, req)

	err = r.DAO.CreateUser(ctx, userDO)
	if err != nil {
		return nil, err
	}
	return utils.ConvertUserDTOToUserVO(req), nil
}

//TODO: temporally only write how to send and process verify code request,
// however, currently not consider about what should be returned in this request.

// maybe, in go, for front-end received data, we should use a stuct whose component var are pointer.
func (r *UserService) UpdateUser(ctx context.Context, req *model.UpdateUserDTO, id int64) (*model.UserVO, error) {
	if req.UserName != nil {
		if *req.UserName == "" {
			return nil, errors.New("username can't be empty")
		}
	}

	if req.Email != nil {
		if *req.Email == "" {
			return nil, errors.New("email can't be empty")
		}
	}

	if req.AvatarURL != nil {
		if *req.AvatarURL == "" {
			return nil, errors.New("AvatarURL can't be empty")
		}
		
	}

	if req.Password != nil {
		if *req.Password == "" {
			return nil, errors.New("Password can't be empty")
		} else {
		}
	}

	userDO, err := r.DAO.GetUserByUserId(ctx, id)
	if err != nil {
		return nil, err
	}
	updateUserDO := model.UpdateUserDO{UpdateUserDTO: *req}
	updateUserMap := utils.StructToMap(updateUserDO, "gorm")
	userId, _ := utils.ReadFromContext[int64](ctx, "userId")
	err = r.DAO.UpdateUser(ctx, updateUserMap, userId)
	if err != nil {
		return nil, err
	}
	userDO, err = r.DAO.GetUserByUserId(ctx, userId)
	userVO := model.UserVO{}
	copier.Copy(userVO, userDO)
	return &userVO, nil
}

// all data received from front-end maybe processed with a struct whose var is pointer, not value
func (r *UserService) UpdatePasswordWithOldPassword(ctx context.Context, req *model.UpdatePasswordWithOldPasswordDTO, id int64) error {
	user, err := r.DAO.GetUserByUserId(ctx, id)
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
	userId, _ := utils.ReadFromContext[int64](ctx, "userId")
	user.Password = utils.Encrypt(req.NewPassword)
	updateUserDO := model.UpdateRepoDO{}
	updateUserDO.Password = utils.LiteralPtr(utils.Encrypt(req.NewPassword))
	updateUserMap := utils.StructToMap(updateUserDO, "gorm")
	err = r.DAO.UpdateUser(ctx, updateUserMap, userId)
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
