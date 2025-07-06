package service

import (
	"go-gcs/appError"
	"go-gcs/dao"
	"go-gcs/model"
	"go-gcs/utils"

	"context"

	"github.com/jinzhu/copier"
)

type RepoService struct {
	RepoDAO *dao.RepoDB
	UserDAO *dao.UserDB
}

// store user and repo infomation may reduce complexity
func (r *RepoService) CreateRepo(ctx context.Context, req *model.CreateRepoDTO, userId int64) (error) {
	exist, _ := r.RepoDAO.IsExists(ctx, req.RepoName, userId)
	if exist {
		return appError.ErrorRepoAlreadyExist
	}
	repoDO := &model.RepoDO{}
	err := copier.Copy(repoDO, req)
	if err != nil {
		return err
	}
	repoDO.User_id = userId
	user, _ := r.UserDAO.GetUserByID(ctx, userId)
	repoDO.Https_url = utils.GenerateHttpURL("localhost", 1234, user.Username, req.RepoName)
	repoDO.Ssh_url = utils.GenerateSshURL("localhost", 1234, user.Username, req.RepoName)
	err = r.RepoDAO.CreateRepo(ctx, repoDO, userId)
	return err
}

func (r *RepoService) UpdateRepo(ctx context.Context, req *model.UpdateRepoDTO, userId int64) (error) {
	exist, _ := r.RepoDAO.IsExists(ctx, req.RepoName, userId)
	if !exist {
		return appError.ErrorRepoNotExist
	}
	err := r.RepoDAO.UpdateRepo(ctx, req, userId)
	return err
}
