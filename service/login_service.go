package service
import(
	"context"

	"go-gcs/dao"
	"go-gcs/model"
	"go-gcs/appError"
	"go-gcs/auth"
	"go-gcs/utils"
	"go.uber.org/zap"
)


type LoginService struct {
	DAO *dao.UserDB
}

func (r *LoginService) LoginVerifyPassword(ctx context.Context, req model.LoginRequestDTO) error {
	user, err := r.DAO.GetUserByUserName(ctx, req.UserName)
	if err != nil {
		zap.L().Error("LoginVerifyPassword() failed", zap.Error(err))
		return err
		// return appError.ErrUserNotFound
	}
	if user.Password != utils.Encrypt(req.PassWord) {
		return appError.ErrorPasswordInvalid
	}
	return nil
}

func (r *LoginService) Login(ctx context.Context, req *model.LoginRequestDTO) (string, error) {
	err := r.LoginVerifyPassword(ctx, *req)
	if err != nil {
		zap.L().Error("Login() failed", zap.Error(err))
		return "", err
	}
	var id int64
	id, err = r.DAO.GetUserIdByUserName(ctx, req.UserName)
	if err != nil {
		zap.L().Error("Login() failed", zap.Error(err))
		return "", err
	}
	token, err := auth.GenerateToken(req.UserName, id)
	return token, err
}
