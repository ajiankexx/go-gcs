package service

import (
	"go-gcs/appError"
	"go-gcs/dao"
	"go-gcs/model"

	"context"
)

type RepoService struct {
	DAO *dao.RepoDB
}

func (r *RepoService) CreateRepo(ctx context.Context, req *model.CreateRepoDTO, userId string) (error) {
	exist, _ := r.DAO.IsExists(ctx, req.RepoName, userId)
	if exist {
		return appError.ErrorRepoAlreadyExist
	}
	err := r.DAO.CreateRepo(ctx, req, userId)
	return err
}

func (r *RepoService) UpdateRepo(ctx context.Context, req *model.UpdateRepoDTO, userId string) (error) {
	exist, _ := r.DAO.IsExists(ctx, req.RepoName, userId)
	if !exist {
		return appError.ErrorRepoNotExist
	}
	err := r.DAO.UpdateRepo(ctx, req, userId)
	return err
}
