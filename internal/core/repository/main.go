package repository

import (
	"context"

	"github.com/yuldoshevgg/menu-mono/internal/core/repository/psql"
	"github.com/yuldoshevgg/menu-mono/internal/core/repository/psql/sqlc"
)

type PostgresStorageI interface {
	sqlc.Querier
}

type Store interface {
	Postgres() PostgresStorageI
}

type storage struct {
	postgres PostgresStorageI
}

func New(ctx context.Context, dsn string) Store {
	return &storage{
		postgres: psql.NewStore(ctx, dsn),
	}
}

func (s *storage) Postgres() PostgresStorageI {
	return s.postgres
}
