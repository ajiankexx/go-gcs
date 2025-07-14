package db
import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PgPool struct {
	pgpool chan *pgx.Conn
	dsn string
}

func NewPgPool(dsn string, maxConns int) (*PgPool, error) {
	pool := make(chan *pgx.Conn, maxConns)
	for _ = range maxConns {
		conn, err := pgx.Connect(context.Background(), dsn)
		if err != nil {
			return nil, err
		}
		pool <- conn
	}
	return &PgPool{pgpool: pool, dsn: dsn}, nil
}

func (r *PgPool) Get() (*pgx.Conn, error) {
	var data *pgx.Conn
	data 
}
