package health

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Ping(ctx context.Context) error {
	if err := r.db.Ping(ctx); err != nil {
		return err
	}
	return nil
}
