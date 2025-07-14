package dao

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-gcs/model"
	"go-gcs/appError"
	"time"
)

type SshDB struct {
	DB *pgxpool.Pool
}

func (r *SshDB) UploadSshKey(ctx context.Context, data *model.SshDO) error {
	query := `insert into public.t_ssh_key(user_id, name, public_key) 
	values ($1, $2, $3)`
	_, err := r.DB.Exec(ctx, query, data.User_id, data.Name, data.Public_key)
	return err
}

func (r *SshDB) UpdateSshKey(ctx context.Context, data *model.SshDO) error {
	query := `
	update public.t_ssh_key
	set user_id = $1,
		name = $2,
		public_key = $3
	where id = $4
	`
	_, err := r.DB.Exec(ctx, query,
		data.User_id,
		data.Name,
		data.Public_key,
		data.Id,
	)
	return err
}

func (r *SshDB) DeleteSshKey(ctx context.Context, data *model.SshDO) error {
	query := `
	update public.t_ssh_key
	set gmt_deleted = $1,
	where id = $2
	`
	_, err := r.DB.Exec(ctx, query, time.Now(), data.Id)
	return err
}

func (r *SshDB) IsExist(ctx context.Context, sshKeyId int64, userId int64) (bool, error) {
	query := `
	select exists (
	select 1
	from t_ssh_key
	where id = $1
		and user_id = $2
	)
	`
	err := r.DB.QueryRow(ctx, query, sshKeyId, userId)
	if err != nil {
		return false, appError.ErrorSshKeyNotExist
	}
	return true, nil
}

func (r *SshDB) GetIdByName(ctx context.Context, name string, userId int64) (int64, error) {
	query := `select id from t_ssh_key where name = $1 and user_id = $2`
	var Id int64
	err := r.DB.QueryRow(ctx, query, name, userId).Scan(&Id)
	if err != nil {
		return -1, err
	}
	return Id, nil
}
