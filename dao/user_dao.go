package dao

import(
	"context"
	"go-gcs/model"
	"github.com/jackc/pgx/v5"
)

type UserDB struct {
	DB *pgx.Conn
}

func (r *UserDB) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO public.t_users(username, email, user_password, avatar_url) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.DB.QueryRow(ctx, query, user.Username, user.Email, user.Password, user.AvatarURL).Scan(&user.ID)
}
