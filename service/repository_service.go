package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"go-gcs/appError"
	"go-gcs/dao"
	"go-gcs/model"
	"go-gcs/utils"
	"go.uber.org/zap"

	"github.com/jinzhu/copier"
)

type RepoService struct {
	Validator     *validator.Validate
	RepoDAO       *dao.RepoDB
	UserDAO       *dao.UserDB
	GitoliteUtils *utils.GitoliteUtils
}

// store user and repo infomation may reduce complexity
func (r *RepoService) CreateRepo(ctx context.Context, req *model.CreateRepoDTO, userId int64) error {
	if err := r.Validator.Struct(req); err != nil {
		utils.LogValidationErrors(err)
		return appError.ErrorWrongFormatRequestData
	}
	exist, _ := r.RepoDAO.RepoExists(ctx, nil, req.RepoName, userId)
	if exist {
		return appError.ErrorRepoAlreadyExist
	}
	repoDO := &model.RepoDO{}
	err := copier.Copy(repoDO, req)
	if err != nil {
		return err
	}
	userName, err := r.UserDAO.GetUserNameById(ctx, nil, userId)
	repoDO.UserId = userId
	user, _ := r.UserDAO.GetUserByUserId(ctx, nil, userId)
	repoDO.HttpsUrl = utils.GenerateHttpURL("localhost", 1234, user.UserName, req.RepoName)
	repoDO.SshUrl = utils.GenerateSshURL("localhost", 1234, user.UserName, req.RepoName)
	err = r.RepoDAO.CreateRepo(ctx, nil, repoDO, userId)

	var repoId int64
	errChan := make(chan struct{}, 1)
	if err == nil {
		repoId, err = r.RepoDAO.GetRepoIdByName(ctx, nil, req.RepoName, userId)
		if err != nil {
			return appError.ErrorCreateFileFailed
		}
		go func() {
			defer func() {
				errChan <- struct{}{}
			}()
			r.GitoliteUtils.CreateRepository(ctx, repoId, req.RepoName, req.IsPrivate, userId, userName)
		}()
	} else {
		close(errChan)
	}
	<-errChan
	close(errChan)
	return nil
}

func (r *RepoService) UpdateRepo(ctx context.Context, req *model.UpdateRepoDTO) error {
	if err := r.Validator.Struct(req); err != nil {
		utils.LogValidationErrors(err)
		return appError.ErrorWrongFormatRequestData
	}
	userId, ok := utils.ReadFromContext[int64](ctx, "userId")
	if !ok {
		zap.L().Error("failed to get user ID", zap.Error(appError.ErrorUserIdNotFound))
		return appError.ErrorUserIdNotFound
	}
	if req.Id == nil {
		return appError.ErrorRepoIDIsEmpty
	}

	updateRepoDO := &model.UpdateRepoDO{UpdateRepoDTO: *req}
	if req.RepoName != nil {
		userDO, _ := r.UserDAO.GetUserByUserId(ctx, nil, userId)
		httpUrl := utils.GenerateHttpURL("localhost", 1234, userDO.UserName, *req.RepoName)
		sshUrl := utils.GenerateSshURL("localhost", 1234, userDO.UserName, *req.RepoName)
		updateRepoDO.HttpsUrl = &httpUrl
		updateRepoDO.SshUrl = &sshUrl
	}
	updateRepoMap := utils.StructToMap(updateRepoDO, "gorm")
	err := r.RepoDAO.UpdateRepo(ctx, nil, updateRepoMap, userId)
	return err
}
