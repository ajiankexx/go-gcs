package model

type Repo struct {
	ID        string `json:"id"`
	RepoName  string `json:"repository_name"`
	RepoDesc  string `json:"repository_description"`
	IsPrivate bool   `json:"is_private"`
}

// embedding struct
// temporarily, for each route and http method, define a struct to receive request data
type UpdateRepo struct {
	Repo
}

type GetRepo struct {
	ID int `json:"id"`
	UserName string `json:"username"`
	RepoName string `json:"repository_name"`
}

type GetRepoByName struct {
	RepoName string `json:"repository_name"`
}
