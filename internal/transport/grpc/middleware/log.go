package middleware

import (
	"context"

	"github.com/yuldoshevgg/menu-mono/internal/core/repository"
	"github.com/yuldoshevgg/menu-mono/internal/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type LogMiddleware struct {
	repo repository.Store
}

func NewLogMiddleware(repo repository.Store) *LogMiddleware {
	return &LogMiddleware{
		repo: repo,
	}
}

func (*LogMiddleware) GrpcLoggerMiddleware(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	mtdt := extractMetadata(ctx, req)

	logger.Log.Info(
		"---> req: ",
		zap.String("Method", info.FullMethod),
		zap.Any("Body", mtdt.Body),
	)

	resp, err = handler(ctx, req)

	if err != nil {
		logger.Log.Error(
			"<---- error: ",
			zap.String("msg", err.Error()),
		)

		return resp, err
	}

	return resp, err
}
