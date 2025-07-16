package utils

import (
	"go-gcs/model"
)

func ConvertUserDTOToUserVO(d *model.UserDTO) (*model.UserVO) {
	return &model.UserVO{
		UserName: d.UserName,
		Email: d.Email,
		AvatarURL: d.AvatarURL,
	}
}
