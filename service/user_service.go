package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"go-gcs/model"
	"go-gcs/dao"
)

type UserService struct {
	DAO *dao.UserDB
}

func (r *UserService) CreateUser(ctx context.Context, req *model.User) (*model.User, error) {
	hash := md5.Sum([]byte(req.Password))
	req.Password = hex.EncodeToString(hash[:])
	err := r.DAO.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *UserService) UpdateUser(ctx context.Context, req *model.User) (*model.User, error) {
	return req, nil
}
