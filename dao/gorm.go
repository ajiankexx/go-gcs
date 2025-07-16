package dao

import (
	"context"
	"go-gcs/model"

	"gorm.io/gorm"
)

type UserDB struct {
	DB *gorm.DB
}

func (r *UserDB) GetUserByUserId(ctx context.Context, userId int64) (*model.UserDO, error) {
	var userDO model.UserDO
	err := r.DB.WithContext(ctx).
		Model(model.UserDO{}).
		Where("id = ?", userId).
		Scan(&userDO).Error
	if err != nil {
		return nil, err
	}
	return &userDO, nil
}

func (r *UserDB) GetUserByUserName(ctx context.Context, userName string) (*model.UserDO, error) {
	var userDO model.UserDO
	err := r.DB.WithContext(ctx).
		Model(model.UserDO{}).
		Where("username = ?", userName).
		Scan(&userDO).Error
	if err != nil {
		return nil, err
	}
	return &userDO, nil
}

func (r *UserDB) GetUserNameById(ctx context.Context, userId int64) (string, error) {
	var userName string
	err := r.DB.WithContext(ctx).
		Model(model.UserDO{}).
		Where("id = ?", userId).
		Pluck("username", userName).Error
	if err != nil {
		return "", err
	}
	return userName, nil
}

func (r *UserDB) GetUserIdByUserName(ctx context.Context, userName string) (int64, error) {
	var userId int64
	err := r.DB.WithContext(ctx).
		Model(model.UserDO{}).
		Where("username = ?", userName).
		Pluck("id", userId).Error
	if err != nil {
		return -1, err
	}
	return userId, nil
}

func (r *UserDB) UserExists(ctx context.Context, userName string) (bool, error) {
	var userDO model.UserDO
	tx := r.DB.WithContext(ctx).
		Model(model.UserDO{}).
		Where("username = ?", userName).
		Scan(&userDO)
	if tx.Error != nil {
		return false, tx.Error
	}
	return tx.RowsAffected > 1, nil
}

func (r *UserDB) CreateUser(ctx context.Context, userDO *model.UserDO) error {
	err := r.DB.WithContext(ctx).
		Model(model.UserDO{}).
		Create(userDO).Error
	return err
}

func (r *UserDB) UpdateUser(ctx context.Context, updateUserMap map[string]interface{},
	userId int64) error {
	err := r.DB.WithContext(ctx).
		Model(model.UserDO{}).
		Where("id = ?", userId).
		Updates(updateUserMap).Error
	return err
}

func (r *UserDB) DeleteUser(ctx context.Context, userId int64) error {
	err := r.DB.WithContext(ctx).
		Model(model.UserDO{}).
		Where("id = ?", userId).
		Delete(model.UserDO{}).Error
	return err
}

type RepoDB struct {
	DB *gorm.DB
}

func (r *RepoDB) RepoExists(ctx context.Context, repoName string, userId int64) (bool, error) {
	var repoDO model.RepoDO
	tx := r.DB.WithContext(ctx).
		Model(model.RepoDO{}).
		Where("repository_name = ? AND user_id = ?", repoName, userId).
		Scan(&repoDO)
	if tx.Error != nil {
		return false, tx.Error
	}
	return tx.RowsAffected > 1, nil
}

func (r *RepoDB) GetRepoIdByName(ctx context.Context, repoName string, userId int64) (int64, error) {
	var repoId int64
	err := r.DB.WithContext(ctx).
		Model(&model.RepoDO{}).
		Where("user_id = ? AND repository_name = ?", userId, repoName).
		Pluck("id", repoId).Error
	return repoId, err
}

func (r *RepoDB) GetRepoNameById(ctx context.Context, repoId int64) (string, error) {
	var res string
	err := r.DB.WithContext(ctx).
		Model(&model.RepoDO{}).
		Where("id = ?", repoId).
		Select("repository_name").
		Scan(&res).Error
	if err != nil {
		return "", err
	}
	return res, nil
}

func (r *RepoDB) CreateRepo(ctx context.Context, repo *model.RepoDO, userId int64) error {
	err := r.DB.WithContext(ctx).
		Model(model.RepoDO{}).
		Create(repo).Error
	return err
}

func (r *RepoDB) UpdateRepo(ctx context.Context, updateRepoMap map[string]interface{}, userId int64) error {
	err := r.DB.WithContext(ctx).
		Model(&model.RepoDO{}).
		Updates(updateRepoMap).Error
	return err
}

func (r *RepoDB) GetRepoByID(ctx context.Context, repoId int64) (*model.RepoDO, error) {
	var repoDO model.RepoDO
	err := r.DB.WithContext(ctx).
		Model(&model.RepoDO{}).
		Where("id = ?", repoId).
		Scan(&repoDO).Error
	if err != nil {
		return nil, err
	}
	return &repoDO, nil
}

type SshKeyDB struct {
	DB *gorm.DB
}

func (r *SshKeyDB) Upload(ctx context.Context, sshKeyDO *model.SshKeyDO) error {
	err := r.DB.WithContext(ctx).
		Model(&model.SshKeyDO{}).
		Create(sshKeyDO).Error
	return err
}

func (r *SshKeyDB) Update(ctx context.Context, updateSshKey map[string]interface{}, sshKeyId int64) error {
	err := r.DB.WithContext(ctx).
		Model(model.SshKeyDO{}).
		Where("id = ?", sshKeyId).Error
	return err
}

func (r *SshKeyDB) SshKeyExists(ctx context.Context, sshKeyName string, userId int64) (bool, error) {
	var sshKeyDO model.SshKeyDO
	tx := r.DB.WithContext(ctx).
		Model(&model.SshKeyDO{}).
		Where("name = ? AND user_id = ?", sshKeyName, userId).
		Scan(&sshKeyDO)
	if tx.Error != nil {
		return false, tx.Error
	}
	return tx.RowsAffected > 1, tx.Error
}

func (r *SshKeyDB) GetSshKeyIdBySshKeyName(ctx context.Context, sshKeyName string, userId int64) (int64, error) {
	var sshKeyId int64
	err := r.DB.WithContext(ctx).
		Model(model.SshKeyDO{}).
		Where("user_id = ? AND name = ?", userId, sshKeyName).
		Pluck("id", sshKeyId).Error
	return sshKeyId, err
}

func (r *SshKeyDB) GetSshKeyNameBySshKeyId(ctx context.Context, sshKeyId int64) (string, error) {
	var sshKeyName string
	err := r.DB.WithContext(ctx).
		Model(model.SshKeyDO{}).
		Where("id = ?", sshKeyId).
		Pluck("name", sshKeyName).Error
	return sshKeyName, err
}

func (r *SshKeyDB) DeleteSshKey(ctx context.Context, sshKeyName string, userId int64) error {
	err := r.DB.WithContext(ctx).
		Model(model.SshKeyDO{}).
		Where("user_id = ? AND name = ?", userId, sshKeyName).
		Delete(model.SshKeyDO{}).Error
	return err
}
