package dao

import (
	"context"

	"go-gcs/appError"
	"go-gcs/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RepoDB struct {
	DB *pgxpool.Pool
}


func (r *RepoDB) IsExists(ctx context.Context, RepoName string, UserId int64) (bool, error) {
	var exists bool
	query := `
		select exists (
			select 1 
			from t_repository
			where repository_name = $1
				and user_id = $2
				and gmt_deleted IS NULL
	)
	`
	// Todo: more type of err may happen
	err := r.DB.QueryRow(ctx, query, RepoName, UserId).Scan(&exists)
	if err != nil {
		return false, appError.ErrorRepoNotExist
	}
	return exists, nil
}

func (r *RepoDB) CreateRepo(ctx context.Context, Repo *model.RepoDO, UserId int64) error {
	exists, _ := r.IsExists(ctx, Repo.RepoName, UserId)
	if exists {
		return appError.ErrorRepoAlreadyExist
	}
	query := `insert into public.t_repository(repository_name, repository_description,
		is_private, user_id, https_url, ssh_url) values ($1, $2, $3, $4, $5, $6) returning id`
	_, err := r.DB.Exec(ctx, query, Repo.RepoName, Repo.RepoDesc, Repo.IsPrivate, Repo.User_id, 
		Repo.Https_url, Repo.Ssh_url)
	return err
}

// temporarily, getRepoDTO contains userid instead of username, convert id to name in service layer.
func (r *RepoDB) GetRepoByName(ctx context.Context, RepoName string, UserId int64) (*model.GetRepoDTO, error) {

	exists, _ := r.IsExists(ctx, RepoName, UserId)
	if !exists {
		return nil, appError.ErrorRepoNotExist
	}
	query := `select repository_name, repository_description, is_private, user_id from public.t_repository where repository_name = $1 and user_id = $2`
	row := r.DB.QueryRow(ctx, query, RepoName, UserId)
	var getRepoDTO model.GetRepoDTO
	err := row.Scan(
		&getRepoDTO.RepoName,
		&getRepoDTO.RepoDesc,
		&getRepoDTO.IsPrivate,
		&getRepoDTO.UserId,
	)
	if err != nil {
		return nil, err
	}
	return &getRepoDTO, nil
}

func (r *RepoDB) UpdateRepo(ctx context.Context, Repo *model.UpdateRepoDTO, UserId int64) (error) {
	query := `update public.t_repository
			  set repository_name = $1,
				  repository_description = $2
				  is_private = $3
				  user_id = $4
			  where id = $5
	`

	_, err := r.DB.Exec(ctx, query, 
		Repo.RepoName,
		Repo.RepoDesc,
		Repo.IsPrivate,
		UserId,
	)
	return err
}
