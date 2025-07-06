package dao

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
	"go-gcs/model"
	"errors"
	"go-gcs/appError"
)

type UserDB struct {
	DB *pgxpool.Pool
}

func (r *UserDB) Create(ctx context.Context, user *model.UserDTO) error {
	query := `INSERT INTO public.t_users(username, email, user_password, avatar_url) VALUES ($1, $2, $3, $4) RETURNING id`
	_, err := r.DB.Exec(ctx, query, user.Username, user.Email, user.Password, user.AvatarURL)
	return err
}

func (r *UserDB) GetUserByName(ctx context.Context, username string) (*model.UserDTO, error) {
	query := `SELECT username, email, user_password, avatar_url FROM public.t_users WHERE username = $1` //BUG: pubulic -> public
	row := r.DB.QueryRow(ctx, query, username)
	var user model.UserDTO
	err := row.Scan(
		&user.Username,
		&user.Email,
		&user.Password,
		&user.AvatarURL,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, appError.ErrorUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserDB) GetUserByID(ctx context.Context, ID int64) (*model.UserDTO, error) {
	query := `SELECT username, email, user_password, avatar_url FROM public.t_users WHERE id = $1`
	row := r.DB.QueryRow(ctx, query, ID)
	var user model.UserDTO
	err := row.Scan(
		&user.Username,
		&user.Email,
		&user.Password,
		&user.AvatarURL,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("ErrNoRows")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserDB) GetUserIDByUserName(ctx context.Context, username string) (*int64, error) {
	query := `SELECT id FROM public.t_users WHERE username = $1`
	row := r.DB.QueryRow(ctx, query, username)
	var id int64
	err := row.Scan(
		&id,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("ErrNoRows")
		}
		return nil, err
	}
	return &id, nil
}

func (r *UserDB) UpdateUser(ctx context.Context, user *model.UserDTO) error {
	query := `
		UPDATE public.t_users
		SET username = $1,
			email = $2,
			user_password = $3,
			avatar_url = $4
		WHERE id = $5
	`

	_, err := r.DB.Exec(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.AvatarURL,
	)
	return err
}
