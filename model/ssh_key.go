package model

type SshDO struct {
	Id         int64  `json:"id"`
	User_id    int64  `json:"user_id"`
	Name       string `json:"name"`
	Public_key string `jso:"public_key"`
}

type SshDTO struct {
	User_id    int64  `json:"user_id"`
	Name       string `json:"name"`
	Public_key string `json:"public_key"`
}

type UpdateSshDTO struct {
	Name       *string `json:"name"`
	Public_key *string `json:"public_key"`
}
