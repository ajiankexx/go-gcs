package model

// 8:55
type RepoDO struct {
	ID        int64  `json:"id"`
	RepoName  string `json:"repository_name"`
	RepoDesc  string `json:"repository_description"`
	IsPrivate bool   `json:"is_private"`
	User_id   int64  `json:"user_id"`
	Star      int    `json:"star"`
	Fork      int    `json:"fork"`
	Watcher   int    `json:"watcher"`
	Https_url string `json:"https_url"`
	Ssh_url   string `json:"ssh_url"`
}

type CreateRepoDTO struct {
	RepoName  string `json:"repository_name"`
	RepoDesc  string `json:"repository_description"`
	IsPrivate bool   `json:"is_private"`
}

// embedding struct
// temporarily, for each route and http method, define a struct to receive request data
type UpdateRepoDTO struct {
	ID        *int64  `json:"id"`
	RepoName  *string `json:"repository_name"`
	RepoDesc  *string `json:"repository_description"`
	IsPrivate *bool   `json:"is_private"`
}

type UpdateRepoDO struct {
	ID        int64  `json:"id"`
	RepoName  string `json:"repository_name"`
	RepoDesc  string `json:"repository_description"`
	IsPrivate bool   `json:"is_private"`
	Https_url string `json:"https_url"`
	Ssh_url   string `json:"ssh_url"`
	User_id   int64  `json:"user_id"`
}

type GetRepoDTO struct {
	UserId    int64  `json:"user_id"`
	RepoName  string `json:"repository_name"`
	RepoDesc  string `json:"repository_description"`
	IsPrivate bool   `json:"is_private"`
}

type GetRepoByNameDTO struct {
	RepoName string `json:"repository_name"`
}
