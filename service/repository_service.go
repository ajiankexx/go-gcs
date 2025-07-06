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

func (r *RepoService) UpdateRepo(ctx context.Context, req *model.UpdateRepoDTO, userId int64) error {
	if req.ID == nil {
		return appError.ErrorRepoIDIsEmpty
	}

	repoDO, err := r.RepoDAO.GetRepoByID(ctx, *req.ID)
	if err != nil {
		return err
	}

	updateRepoDO := &model.UpdateRepoDO{
		ID:        repoDO.ID,
		RepoName:  repoDO.RepoName,
		RepoDesc:  repoDO.RepoDesc,
		IsPrivate: repoDO.IsPrivate,
		Https_url: repoDO.Https_url,
		Ssh_url:   repoDO.Ssh_url,
		User_id:   repoDO.User_id,
	}

	userDTO, err := r.UserDAO.GetUserByID(ctx, userId)
	if req.RepoName != nil {
		updateRepoDO.RepoName = *req.RepoName
		updateRepoDO.Https_url = utils.GenerateHttpURL("localhost", 1234, userDTO.Username, *req.RepoName)
		updateRepoDO.Ssh_url = utils.GenerateSshURL("localhost", 1234, userDTO.Username, *req.RepoName)
	}
	if req.RepoDesc != nil {
		updateRepoDO.RepoDesc = *req.RepoDesc
	}
	if req.IsPrivate != nil {
		updateRepoDO.IsPrivate = *req.IsPrivate
	}

	err = r.RepoDAO.UpdateRepo(ctx, updateRepoDO)
	return err
}
