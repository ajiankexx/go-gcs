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


func (r *RepoDB) IsExists(ctx context.Context, RepoName string, UserId string) (bool, error) {
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
	err := r.DB.QueryRow(ctx, query, RepoName, UserId)
	if err != nil {
		return false, appError.ErrorRepoNotExist
	}
	return exists, nil
}

func (r *RepoDB) CreateRepo(ctx context.Context, Repo *model.Repo, UserId string) error {
	exists, _ := r.IsExists(ctx, Repo.RepoName, UserId)
	if exists {
		return appError.ErrorRepoAlreadyExist
	}
	query := `insert into public.t_repository(repository_name, repository_description,
		is_private, user_id) values ($1, $2, $3, $4) returning id`
	return r.DB.QueryRow(ctx, query, Repo.RepoName, Repo.RepoDesc, Repo.IsPrivate, UserId).
		Scan(&Repo.ID)
}

func (r *RepoDB) GetRepoByName(ctx context.Context, RepoName string, UserId string) (*model.Repo, error) {

	exists, _ := r.IsExists(ctx, RepoName, UserId)
	if !exists {
		return nil, appError.ErrorRepoNotExist
	}
	query := `select id, repository_name, repository_description, is_private, user_id from public.t_repository where repository_name = $1 and user_id = $2`
	row := r.DB.QueryRow(ctx, query, RepoName, UserId)
	var RepoInfo model.Repo
	err := row.Scan(
		&RepoInfo.ID,
		&RepoInfo.RepoName,
		&RepoInfo.RepoDesc,
		&RepoInfo.IsPrivate,
	)
	if err != nil {
		return nil, err
	}
	return &RepoInfo, nil
}

func (r *RepoDB) UpdateRepo(ctx context.Context, Repo *model.UpdateRepo, UserId string) (error) {
	query := `update public.t_repository
			  set repository_name $1,
				  repository_description $2
				  is_private $3
				  user_id $4
	`

	_, err := r.DB.Exec(ctx, query, 
		Repo.RepoName,
		Repo.RepoDesc,
		Repo.IsPrivate,
		UserId,
	)
	return err
}
