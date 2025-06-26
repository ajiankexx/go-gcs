package service
import(
	"context"

	"go-gcs/dao"
	"go-gcs/model"
	"go-gcs/appError"
	"go-gcs/auth"
	"go-gcs/utils"
)


type LoginService struct {
	DAO *dao.UserDB
}

func (r *LoginService) LoginVerifyPassword(ctx context.Context, req model.LoginRequest) error {
	user, err := r.DAO.GetUserByName(ctx, req.UserName)
	if err != nil {
		return err
		// return appError.ErrUserNotFound
	}
	if user.Password != utils.Encrypt(req.PassWord) {
		return appError.ErrPasswordInvalid
	}
	return nil
}

func (r *LoginService) Login(ctx context.Context, req *model.LoginRequest) (string, error) {
	err := r.LoginVerifyPassword(ctx, *req)
	if err != nil {
		return "", err
	}
	token, err := auth.GenerateToken(req.UserName)
	return token, err
}


