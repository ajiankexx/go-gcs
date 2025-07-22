package model

import (
	"go-gcs/constants"
)

var ssh_key_table_name = constants.SSH_KEY_TABLE_NAME

type SshKeyDO struct {
	sshKeyTableName
	Id        int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    int64  `gorm:"column:user_id" json:"user_id"`
	Name      string `gorm:"column:name" json:"name"`
	PublicKey string `gorm:"column:public_key" jso:"public_key"`
}

type SshKeyDTO struct {
	sshKeyTableName
	UserId    int64  `gorm:"column:user_id" json:"user_id" validate:"required"`
	Name      string `gorm:"column:name" json:"name" validate:"required,min=1,max=100"`
	PublicKey string `gorm:"column:public_key" json:"public_key" validate:"required,startswith=ssh-"`
}

type UpdateSshKeyDTO struct {
	sshKeyTableName
	OldName      *string `json:"old_name" validate:"omitempty,min=1,max=100"`
	NewName      *string `json:"new_name" validate:"omitempty,min=1,max=100"`
	OldPublicKey *string `json:"old_public_key" validate:"omitempty,startswith=ssh-"`
	NewPublicKey *string `json:"new_public_key" validate:"omitempty,startswith=ssh-"`
}
