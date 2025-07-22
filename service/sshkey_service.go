package service
import (
	"go-gcs/appError"
	"go-gcs/dao"
	"go-gcs/model"
	"go-gcs/utils"
	"go.uber.org/zap"
	"github.com/go-playground/validator/v10"

	"context"

	"github.com/jinzhu/copier"
)

type SshKeyService struct{
	Validator *validator.Validate
	SshKeyDAO *dao.SshKeyDB
}

func (r *SshKeyService) UploadSshKey(ctx context.Context, req *model.SshKeyDTO, userId int64) error {
	if err := r.Validator.Struct(req); err != nil {
		utils.LogValidationErrors(err)
		return appError.ErrorWrongFormatRequestData
	}
	exist, err := r.SshKeyDAO.SshKeyExists(ctx, nil, req.Name, req.UserId)
	if err != nil {
		return err
	}
	if exist {
		return appError.ErrorSshKeyAlreadyExist
	}
	var sshKeyDO = &model.SshKeyDO{}
	copier.Copy(sshKeyDO, req)
	err = r.SshKeyDAO.Upload(ctx, nil, sshKeyDO)
	if err != nil {
		return err
	}
	return nil
}

func (r *SshKeyService) UpdateSshKey(ctx context.Context, req *model.UpdateSshKeyDTO) error {
	if err := r.Validator.Struct(req); err != nil {
		utils.LogValidationErrors(err)
		return appError.ErrorWrongFormatRequestData
	}
	userId, _ := utils.ReadFromContext[int64](ctx, "userId")
	var exists bool
	exists, _ = r.SshKeyDAO.SshKeyExists(ctx, nil, *req.OldName, userId)
	if !exists {
		return appError.ErrorSshKeyNotExist
	}
	sshKeyId, err := r.SshKeyDAO.GetSshKeyIdBySshKeyName(ctx, nil, *req.OldName, userId)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	var updateSshDO model.SshKeyDO
	updateSshDO.Name = *req.NewName
	updateSshDO.PublicKey = *req.NewPublicKey
	updateSshKeyMap := utils.StructToMap(updateSshDO, "gorm")
	err = r.SshKeyDAO.Update(ctx, nil, updateSshKeyMap, sshKeyId)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil
}

func (r *SshKeyService) DeleteSshKey(ctx context.Context, req *model.SshKeyDTO) error {
	if err := r.Validator.Struct(req); err != nil {
		utils.LogValidationErrors(err)
		return appError.ErrorWrongFormatRequestData
	}
	userId, _ := utils.ReadFromContext[int64](ctx, "userId")
	sshKeyName := req.Name
	exists, _ := r.SshKeyDAO.SshKeyExists(ctx, nil, sshKeyName, userId)
	if !exists {
		return appError.ErrorSshKeyNotExist
	}
	err := r.SshKeyDAO.DeleteSshKey(ctx, nil, sshKeyName, userId)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil
}
