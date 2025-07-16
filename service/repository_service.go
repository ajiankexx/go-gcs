package service

import (
	"go-gcs/appError"
	"go-gcs/dao"
	"go-gcs/model"
	"go-gcs/utils"
	"go.uber.org/zap"

	"context"

	"github.com/jinzhu/copier"
)

type RepoService struct {
	RepoDAO *dao.RepoDB
	UserDAO *dao.UserDB
}

// store user and repo infomation may reduce complexity
func (r *RepoService) CreateRepo(ctx context.Context, req *model.CreateRepoDTO, userId int64) (error) {
	exist, _ := r.RepoDAO.RepoExists(ctx, req.RepoName, userId)
	if exist {
		return appError.ErrorRepoAlreadyExist
	}
	repoDO := &model.RepoDO{}
	err := copier.Copy(repoDO, req)
	if err != nil {
		return err
	}
	repoDO.UserId = userId
	user, _ := r.UserDAO.GetUserByUserId(ctx, userId)
	repoDO.HttpsUrl = utils.GenerateHttpURL("localhost", 1234, user.UserName, req.RepoName)
	repoDO.SshUrl = utils.GenerateSshURL("localhost", 1234, user.UserName, req.RepoName)
	err = r.RepoDAO.CreateRepo(ctx, repoDO, userId)
	return nil
}

func (r *RepoService) UpdateRepo(ctx context.Context, req *model.UpdateRepoDTO) error {
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
		userDO, _ := r.UserDAO.GetUserByUserId(ctx, userId)
		httpUrl := utils.GenerateHttpURL("localhost", 1234, userDO.UserName, *req.RepoName)
		sshUrl := utils.GenerateSshURL("localhost", 1234, userDO.UserName, *req.RepoName)
		updateRepoDO.HttpsUrl = &httpUrl
		updateRepoDO.SshUrl = &sshUrl
	}
	updateRepoMap := utils.StructToMap(updateRepoDO, "gorm")
	err := r.RepoDAO.UpdateRepo(ctx, updateRepoMap, userId)
	return err
}
