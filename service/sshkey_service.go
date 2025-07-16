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

type SshKeyService struct{
	SshKeyDAO *dao.SshKeyDB
}

func (r *SshKeyService) UploadSshKey(ctx context.Context, req *model.SshKeyDTO, userId int64) error {
	exist, err := r.SshKeyDAO.SshKeyExists(ctx, req.Name, req.UserId)
	if err != nil {
		return err
	}
	if exist {
		return appError.ErrorSshKeyAlreadyExist
	}
	var sshKeyDO = &model.SshKeyDO{}
	copier.Copy(sshKeyDO, req)
	err = r.SshKeyDAO.Upload(ctx, sshKeyDO)
	if err != nil {
		return err
	}
	return nil
}

func (r *SshKeyService) UpdateSshKey(ctx context.Context, req *model.UpdateSshKeyDTO) error {
	userId, _ := utils.ReadFromContext[int64](ctx, "userId")
	var exists bool
	exists, _ = r.SshKeyDAO.SshKeyExists(ctx, *req.OldName, userId)
	if !exists {
		return appError.ErrorSshKeyNotExist
	}
	sshKeyId, err := r.SshKeyDAO.GetSshKeyIdBySshKeyName(ctx, *req.OldName, userId)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	var updateSshDO model.SshKeyDO
	updateSshDO.Name = *req.NewName
	updateSshDO.PublicKey = *req.NewPublicKey
	updateSshKeyMap := utils.StructToMap(updateSshDO, "gorm")
	err = r.SshKeyDAO.Update(ctx, updateSshKeyMap, sshKeyId)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil
}

func (r *SshKeyService) DeleteSshKey(ctx context.Context, req *model.SshKeyDTO) error {
	userId, _ := utils.ReadFromContext[int64](ctx, "userId")
	sshKeyName := req.Name
	exists, _ := r.SshKeyDAO.SshKeyExists(ctx, sshKeyName, userId)
	if !exists {
		return appError.ErrorSshKeyNotExist
	}
	err := r.SshKeyDAO.DeleteSshKey(ctx, sshKeyName, userId)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil
}
