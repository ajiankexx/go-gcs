package model
import (
)

type RepoDO struct {
	repoTableName
	Id        int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RepoName  string `gorm:"column:repository_name" json:"repository_name"`
	RepoDesc  string `gorm:"column:repository_description" json:"repository_description"`
	IsPrivate bool   `gorm:"column:is_private" json:"is_private"`
	UserId    int64  `gorm:"column:user_id" json:"user_id"`
	Star      int    `gorm:"column:star" json:"star"`
	Fork      int    `gorm:"column:fork" json:"fork"`
	Watcher   int    `gorm:"column:watcher" json:"watcher"`
	HttpsUrl  string `gorm:"column:https_url" json:"https_url"`
	SshUrl    string `gorm:"column:ssh_url" json:"ssh_url"`
}

type CreateRepoDTO struct {
	repoTableName
	RepoName  string `gorm:"column:repository_name" json:"repository_name"`
	RepoDesc  string `gorm:"column:repository_description" json:"repository_description"`
	IsPrivate bool   `gorm:"column:is_private" json:"is_private"`
}

// embedding struct
// temporarily, for each route and http method, define a struct to receive request data
type UpdateRepoDTO struct {
	repoTableName
	Id        *int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RepoName  *string `gorm:"column:repository_name" json:"repository_name"`
	RepoDesc  *string `gorm:"column:repository_description" json:"repository_description"`
	IsPrivate *bool   `gorm:"column:is_private" json:"is_private"`
}

type UpdateRepoDO struct {
	repoTableName
	UpdateRepoDTO
	HttpsUrl *string `gorm:"column:https_url" json:"https_url"`
	SshUrl   *string `gorm:"column:ssh_url" json:"ssh_url"`
}
