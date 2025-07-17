package dao

import (
	"context"
	"errors"
	"go-gcs/model"
	"runtime/debug"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserDB struct {
	DB *gorm.DB
}

func (r *UserDB) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		zap.L().Error("try to get transaction connection failed")
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			zap.L().Error("recover from panic",
				zap.Any("panic", r),
				zap.ByteString("stack", debug.Stack()),
			)
			tx.Rollback()
			zap.L().Info("rollback successful")
			panic(r)
		}
	}()
	if err := fn(tx); err != nil {
		zap.L().Error("transaction execute failed", zap.Error(err))
		tx.Rollback()
		zap.L().Info("rollback successful")
		return err
	}
	if err := tx.Commit().Error; err != nil {
		zap.L().Error("commit transaction failed", zap.Error(err))
		return err
	}
	return nil
}

func (r *UserDB) GetUserByUserId(ctx context.Context, tx *gorm.DB, userId int64) (*model.UserDO, error) {
	db := getDB(tx, r.DB)
	var userDO model.UserDO
	err := db.WithContext(ctx).
		Model(model.UserDO{}).
		Where("id = ?", userId).
		Scan(&userDO).Error
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	return &userDO, nil
}

func (r *UserDB) GetUserByUserName(ctx context.Context, tx *gorm.DB, userName string) (*model.UserDO, error) {
	db := getDB(tx, r.DB)
	var userDO model.UserDO
	err := db.WithContext(ctx).
		Model(model.UserDO{}).
		Where("username = ?", userName).
		Scan(&userDO).Error
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	return &userDO, nil
}

func (r *UserDB) GetUserNameById(ctx context.Context, tx *gorm.DB, userId int64) (string, error) {
	db := getDB(tx, r.DB)
	var userName string
	err := db.WithContext(ctx).
		Model(model.UserDO{}).
		Where("id = ?", userId).
		Pluck("username", &userName).Error
	if err != nil {
		zap.L().Error(err.Error())
		return "", err
	}
	return userName, nil
}

func (r *UserDB) GetUserIdByUserName(ctx context.Context, tx *gorm.DB, userName string) (int64, error) {
	db := getDB(tx, r.DB)
	var userId int64
	err := db.WithContext(ctx).
		Model(model.UserDO{}).
		Where("username = ?", userName).
		Pluck("id", &userId).Error
	if err != nil {
		zap.L().Error(err.Error())
		return -1, err
	}
	return userId, nil
}

func (r *UserDB) UserExists(ctx context.Context, tx *gorm.DB, userName string) (bool, error) {
	db := getDB(tx, r.DB)
	var user model.UserDO
	err := db.WithContext(ctx).
		Model(model.UserDO{}).
		Where("username = ?", userName).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		zap.L().Error("failed to check user existence", zap.Error(err))
		return false, err
	}
	return true, nil
}

// func (r *UserDB) UserExists(ctx context.Context, userName string) (bool, error) {
// 	var userDO model.UserDO
// 	tx := r.DB.WithContext(ctx).
// 		Model(model.UserDO{}).
// 		Where("username = ?", userName).
// 		Scan(&userDO)
// 	if tx.Error != nil {
// 		zap.L().Error(tx.Error.Error())
// 		return false, tx.Error
// 	}
// 	return tx.RowsAffected > 1, nil
// }

func (r *UserDB) CreateUser(ctx context.Context, tx *gorm.DB, userDO *model.UserDO) error {
	db := getDB(tx, r.DB)
	err := db.WithContext(ctx).
		Model(model.UserDO{}).
		Create(userDO).Error
	if err != nil {
		zap.L().Error(err.Error())
	}
	return err
}

func (r *UserDB) UpdateUser(ctx context.Context, tx *gorm.DB, updateUserMap map[string]interface{}, userId int64) error {
	db := getDB(tx, r.DB)
	err := db.WithContext(ctx).
		Model(model.UserDO{}).
		Where("id = ?", userId).
		Updates(updateUserMap).Error
	if err != nil {
		zap.L().Error(err.Error())
	}
	return err
}

func (r *UserDB) DeleteUser(ctx context.Context, tx *gorm.DB, userId int64) error {
	db := getDB(tx, r.DB)
	err := db.WithContext(ctx).
		Model(model.UserDO{}).
		Where("id = ?", userId).
		Delete(model.UserDO{}).Error
	if err != nil {
		zap.L().Error(err.Error())
	}
	return err
}

type RepoDB struct {
	DB *gorm.DB
}

func (r *RepoDB) WithTransaction(ctx context.Context, fn func(*gorm.DB) error) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		zap.L().Error("Try to get transaction connection failed", zap.Error(tx.Error))
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			zap.L().Error("recover from panic",
				zap.Any("panic", r),
				zap.ByteString("stack", debug.Stack()),
			)
			tx.Rollback()
			zap.L().Info("rollback successful")
			panic(r)
		}
	}()
	if err := fn(tx); err != nil {
		zap.L().Error("transaction execute failed, rollback successful", zap.Error(err))
		tx.Rollback()
		zap.L().Info("rollback successful")
		return err
	}
	if err := tx.Commit().Error; err != nil {
		zap.L().Error("commit transaction failed", zap.Error(err))
		return err
	}
	return nil
}

func (r *RepoDB) RepoExists(ctx context.Context, tx *gorm.DB, repoName string, userId int64) (bool, error) {
	db := getDB(tx, r.DB)
	var repoDO model.RepoDO
	err := db.WithContext(ctx).
		Model(model.RepoDO{}).
		Where("repository_name = ? AND user_id = ?", repoName, userId).
		First(&repoDO).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		zap.L().Error("failed to check repository existence", zap.Error(err))
		return false, err
	}
	return true, nil
}

// func (r *RepoDB) RepoExists(ctx context.Context, repoName string, userId int64) (bool, error) {
// 	var repoDO model.RepoDO
// 	tx := r.DB.WithContext(ctx).
// 		Model(model.RepoDO{}).
// 		Where("repository_name = ? AND user_id = ?", repoName, userId).
// 		Scan(&repoDO)
// 	if tx.Error != nil {
// 		zap.L().Error(tx.Error.Error())
// 		return false, tx.Error
// 	}
// 	return tx.RowsAffected > 1, nil
// }

func (r *RepoDB) GetRepoIdByName(ctx context.Context, tx *gorm.DB, repoName string, userId int64) (int64, error) {
	db := getDB(tx, r.DB)
	var repoId int64
	err := db.WithContext(ctx).
		Model(&model.RepoDO{}).
		Where("user_id = ? AND repository_name = ?", userId, repoName).
		Pluck("id", &repoId).Error
	if err != nil {
		zap.L().Error(err.Error())
	}
	return repoId, err
}

func (r *RepoDB) GetRepoNameById(ctx context.Context, tx *gorm.DB, repoId int64) (string, error) {
	db := getDB(tx, r.DB)
	var res string
	err := db.WithContext(ctx).
		Model(&model.RepoDO{}).
		Where("id = ?", repoId).
		Select("repository_name").
		Scan(&res).Error
	if err != nil {
		zap.L().Error(err.Error())
		return "", err
	}
	return res, nil
}

func (r *RepoDB) CreateRepo(ctx context.Context, tx *gorm.DB, repo *model.RepoDO, userId int64) error {
	db := getDB(tx, r.DB)
	err := db.WithContext(ctx).
		Model(model.RepoDO{}).
		Create(repo).Error
	if err != nil {
		zap.L().Error(err.Error())
	}
	return err
}

func (r *RepoDB) UpdateRepo(ctx context.Context, tx *gorm.DB, updateRepoMap map[string]interface{}, userId int64) error {
	db := getDB(tx, r.DB)
	err := db.WithContext(ctx).
		Model(&model.RepoDO{}).
		Updates(updateRepoMap).Error
	return err
}

func (r *RepoDB) GetRepoByID(ctx context.Context, tx *gorm.DB, repoId int64) (*model.RepoDO, error) {
	db := getDB(tx, r.DB)
	var repoDO model.RepoDO
	err := db.WithContext(ctx).
		Model(&model.RepoDO{}).
		Where("id = ?", repoId).
		Scan(&repoDO).Error
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	return &repoDO, nil
}

type SshKeyDB struct {
	DB *gorm.DB
}

func (r *SshKeyDB) Upload(ctx context.Context, tx *gorm.DB, sshKeyDO *model.SshKeyDO) error {
	db := getDB(tx, r.DB)
	err := db.WithContext(ctx).
		Model(&model.SshKeyDO{}).
		Create(sshKeyDO).Error
	if err != nil {
		zap.L().Error(err.Error())
	}
	return err
}

func (r *SshKeyDB) Update(ctx context.Context, tx *gorm.DB, updateSshKey map[string]interface{}, sshKeyId int64) error {
	db := getDB(tx, r.DB)
	err := db.WithContext(ctx).
		Model(model.SshKeyDO{}).
		Where("id = ?", sshKeyId).Error
	if err != nil {
		zap.L().Error(err.Error())
	}
	return err
}

func (r *SshKeyDB) WithTransaction(ctx, fn func(*gorm.DB) error) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		zap.L().Error("try to get transaction connection failed", zap.Error(tx.Error))
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			zap.L().Error("recover from panic",
				zap.Any("panic", r),
				zap.ByteString("stack", debug.Stack()),
			)
			tx.Rollback()
			zap.L().Info("rollback successful")
			panic(r)
		}
	}()
	if err := fn(tx); err != nil {
		zap.L().Error("transaction execute failed", zap.Error(err))
		tx.Rollback()
		zap.L().Info("rollback successful")
		return err
	}
	if err := tx.Commit().Error; err != nil {
		zap.L().Error("commit transaction failed", zap.Error(err))
		tx.Rollback()
		zap.L().Info("rollback successful")
		return err
	}
	return nil
}

func (r *SshKeyDB) SshKeyExists(ctx context.Context, tx *gorm.DB, sshKeyName string, userId int64) (bool, error) {
	db := getDB(tx, r.DB)
	var sshKeyDO model.SshKeyDO
	err := db.WithContext(ctx).
		Model(model.SshKeyDO{}).
		Where("name = ? AND user_id = ?", sshKeyName, userId).
		First(&sshKeyDO).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		zap.L().Error("failed to check user existence", zap.Error(err))
		return false, err
	}
	return true, nil
}

// func (r *SshKeyDB) SshKeyExists(ctx context.Context, sshKeyName string, userId int64) (bool, error) {
// 	var sshKeyDO model.SshKeyDO
// 	tx := r.DB.WithContext(ctx).
// 		Model(&model.SshKeyDO{}).
// 		Where("name = ? AND user_id = ?", sshKeyName, userId).
// 		Scan(&sshKeyDO)
// 	if tx.Error != nil {
// 		zap.L().Error(tx.Error.Error())
// 		return false, tx.Error
// 	}
// 	return tx.RowsAffected > 1, tx.Error
// }

func (r *SshKeyDB) GetSshKeyIdBySshKeyName(ctx context.Context, tx *gorm.DB, sshKeyName string, userId int64) (int64, error) {
	db := getDB(tx, r.DB)
	var sshKeyId int64
	err := db.WithContext(ctx).
		Model(model.SshKeyDO{}).
		Where("user_id = ? AND name = ?", userId, sshKeyName).
		Pluck("id", &sshKeyId).Error
	if err != nil {
		zap.L().Error(err.Error())
	}
	return sshKeyId, err
}

func (r *SshKeyDB) GetSshKeyNameBySshKeyId(ctx context.Context, tx *gorm.DB, sshKeyId int64) (string, error) {
	db := getDB(tx, r.DB)
	var sshKeyName string
	err := db.WithContext(ctx).
		Model(model.SshKeyDO{}).
		Where("id = ?", sshKeyId).
		Pluck("name", &sshKeyName).Error
	if err != nil {
		zap.L().Error(err.Error())
	}
	return sshKeyName, err
}

func (r *SshKeyDB) DeleteSshKey(ctx context.Context, tx *gorm.DB, sshKeyName string, userId int64) error {
	db := getDB(tx, r.DB)
	err := db.WithContext(ctx).
		Model(model.SshKeyDO{}).
		Where("user_id = ? AND name = ?", userId, sshKeyName).
		Delete(model.SshKeyDO{}).Error
	if err != nil {
		zap.L().Error(err.Error())
	}
	return err
}

func getDB(tx *gorm.DB, fallback *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return fallback
}
