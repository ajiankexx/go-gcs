package dao

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go-gcs/model"
	"errors"
)

type UserDB struct {
	DB *pgx.Conn
}

func (r *UserDB) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO public.t_users(username, email, user_password, avatar_url) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.DB.QueryRow(ctx, query, user.Username, user.Email, user.Password, user.AvatarURL).Scan(&user.ID)
}

func (r *UserDB) GetUserByName(ctx context.Context, username string) (*model.User, error) {
	query := `SELECT id, username, email, user_password, avatar_url FROM public.t_users WHERE username = $1` //BUG: pubulic -> public
	row := r.DB.QueryRow(ctx, query, username)
	var user model.User
	err := row.Scan(
		&user.ID,
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
