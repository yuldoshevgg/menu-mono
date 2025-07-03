package grpc

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/yuldoshevgg/menu-mono/generated/auth_service"
	"github.com/yuldoshevgg/menu-mono/generated/general_service"
	"github.com/yuldoshevgg/menu-mono/generated/menu_service"
	"github.com/yuldoshevgg/menu-mono/internal/config"
	"github.com/yuldoshevgg/menu-mono/internal/core/repository"
	"github.com/yuldoshevgg/menu-mono/internal/core/service"

	"github.com/yuldoshevgg/menu-mono/internal/transport/grpc/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	MaxRecvMsgSize = 10 * 1024 * 1024
)

type GrpcTransport struct {
	Cfg  *config.Config
	Repo repository.Store
}

func New(grpcTransport *GrpcTransport) (grpcServer *grpc.Server) {
	logMiddleware := middleware.NewLogMiddleware(grpcTransport.Repo)
	grpcServer = grpc.NewServer(
		grpc.MaxRecvMsgSize(MaxRecvMsgSize),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				logMiddleware.GrpcLoggerMiddleware,
				auth.UnaryServerInterceptor(middleware.AuthFn),
				grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(middleware.GrpcGateWayRecovery)),
			),
		),
	)

	reflection.Register(grpcServer)
	auth_service.RegisterRoleServiceServer(grpcServer, service.NewRoleService(grpcTransport.Cfg, grpcTransport.Repo))
	general_service.RegisterVenueServiceServer(grpcServer, service.NewVenueService(grpcTransport.Cfg, grpcTransport.Repo))
	menu_service.RegisterMenuServiceServer(grpcServer, service.NewMenuService(grpcTransport.Cfg, grpcTransport.Repo))

	return grpcServer
}
