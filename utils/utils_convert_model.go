package utils

import (
	"go-gcs/model"
)

func ConvertUserDTOToUserVO(d *model.UserDTO) (*model.UserVO) {
	return &model.UserVO{
		Username: d.Username,
		Email: d.Email,
		AvatarURL: d.AvatarURL,
	}
}
