package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/yuldoshevgg/menu-mono/generated/auth_service"
	"github.com/yuldoshevgg/menu-mono/generated/general_service"
	"github.com/yuldoshevgg/menu-mono/generated/menu_service"
	"github.com/yuldoshevgg/menu-mono/internal/config"
	"github.com/yuldoshevgg/menu-mono/internal/pkg/logger"
	"github.com/yuldoshevgg/menu-mono/pkg/wrapper"
	"go.uber.org/zap"
	mainGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func registerHandlers(ctx context.Context, gwMux *runtime.ServeMux, cfg *config.Config) {
	connGrpcService, err := mainGrpc.NewClient(
		makeHost(cfg.Grpc.Host, cfg.Grpc.Port),
		mainGrpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Log.Error("error in method Dial: ", zap.Error(err))
	}

	if err := auth_service.RegisterUserServiceHandler(ctx, gwMux, connGrpcService); err != nil {
		logger.Log.Error("error in method RegisterUserServiceHandler: ", zap.Error(err))
	}

	if err := auth_service.RegisterRoleServiceHandler(ctx, gwMux, connGrpcService); err != nil {
		logger.Log.Error("error in method RegisterRoleServiceHandler: ", zap.Error(err))
	}

	if err := general_service.RegisterVenueServiceHandler(ctx, gwMux, connGrpcService); err != nil {
		logger.Log.Error("error in method RegisterVenueServiceHandler: ", zap.Error(err))
	}

	if err := menu_service.RegisterCategoryServiceHandler(ctx, gwMux, connGrpcService); err != nil {
		logger.Log.Error("error in method RegisterCategoryServiceHandler: ", zap.Error(err))
	}

	if err := menu_service.RegisterMenuServiceHandler(ctx, gwMux, connGrpcService); err != nil {
		logger.Log.Error("error in method RegisterMenuServiceHandler: ", zap.Error(err))
	}
}

func New(ctx context.Context, cfg *config.Config) *runtime.ServeMux {
	gwMux := runtime.NewServeMux(
		runtime.WithMetadata(func(_ context.Context, req *http.Request) metadata.MD {
			return metadata.New(map[string]string{
				config.GrpcGatewayHTTPPath:   req.URL.Path,
				config.GrpcGatewayHTTPMethod: req.Method,
			})
		}),
		runtime.WithIncomingHeaderMatcher(wrapper.CustomMatcher),
	)

	registerHandlers(ctx, gwMux, cfg)

	return gwMux
}

func makeHost(host string, port int) string {
	return host + ":" + fmt.Sprintf("%d", port)
}
