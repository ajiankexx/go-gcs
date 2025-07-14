package service
import (
	"go-gcs/appError"
	"go-gcs/dao"
	"go-gcs/model"
	"go-gcs/utils"

	"context"

	"github.com/jinzhu/copier"
)

type SshKeyService struct{
	SshDAO *dao.SshDB
}

func (r *SshKeyService) UploadSshKey(ctx context.Context, req *model.SshDTO, userId int64) (error) {
	id, err := r.SshDAO.GetIdByName(ctx, req.Name, req.User_id)
	if err != nil {
		return err
	}
	exist, err := r.SshDAO.IsExist(ctx, id, req.User_id)
	if err != nil {
		return err
	}
	if exist {
		return appError.ErrorSshKeyAlreadyExist
	}
	var reqDO = &model.SshDO{}
	copier.Copy(reqDO, req)
	err = r.SshDAO.UploadSshKey(ctx, reqDO)
	if err != nil {
		return err
	}
	return nil
}
