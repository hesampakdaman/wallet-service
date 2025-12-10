package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Starter struct {
	pool *pgxpool.Pool
}

func NewStarter(p *pgxpool.Pool) *Starter {
	return &Starter{pool: p}
}

func (s Starter) Repository(ctx context.Context, opts pgx.TxOptions) (Repository, error) {
	tx, err := s.pool.BeginTx(ctx, opts)
	if err != nil {
		return Repository{}, err
	}
	return Repository{tx: tx}, nil
}
