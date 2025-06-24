package handlers

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/yuldoshevgg/menu-mono/internal/config"
	"google.golang.org/grpc/metadata"
)

func registerHandlers(ctx context.Context, gwMux *runtime.ServeMux, cfg *config.Config) {

}

func New(ctx context.Context, cfg *config.Config) *runtime.ServeMux {
	gwMux := runtime.NewServeMux(
		runtime.WithMetadata(func(_ context.Context, req *http.Request) metadata.MD {
			return metadata.New(map[string]string{
				config.GrpcGatewayHTTPPath:   req.URL.Path,
				config.GrpcGatewayHTTPMethod: req.Method,
			})
		}),
	)

	registerHandlers(ctx, gwMux, cfg)

	return gwMux
}
