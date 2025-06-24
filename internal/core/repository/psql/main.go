package psql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yuldoshevgg/menu-mono/internal/core/repository/psql/sqlc"
	"github.com/yuldoshevgg/menu-mono/internal/pkg/logger"
	"go.uber.org/zap"
)

type SQLStore struct {
	DB *pgxpool.Pool
	*sqlc.Queries
}

func NewStore(ctx context.Context, psqlURI string) *SQLStore {
	logger.Log.Info("connecting to psql...")

	cfg, err := pgxpool.ParseConfig(psqlURI)
	if err != nil {
		logger.Log.Fatal("failed to parse config", zap.Error(err))
	}

	dbConn, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		logger.Log.Fatal("failed to get db conn", zap.Error(err))
	}

	if err := dbConn.Ping(ctx); err != nil {
		logger.Log.Fatal("failed to ping psql", zap.Error(err))
	}

	return &SQLStore{
		DB:      dbConn,
		Queries: sqlc.New(dbConn),
	}
}
