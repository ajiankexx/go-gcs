package model
import (
	"go-gcs/constants"
)

type tableName interface{
	TableName() string
}

type repoTableName struct {
}
func (r repoTableName) TableName() string {
	return constants.REPOSITORY_TABLE_NAME
}

type userTableName struct {
}
func (r userTableName) TableName() string {
	return constants.USER_TABLE_NAME
}

type sshKeyTableName struct {
}
func (r sshKeyTableName) TableName() string {
	return constants.SSH_KEY_TABLE_NAME
}
