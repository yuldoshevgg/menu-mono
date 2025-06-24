package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yuldoshevgg/menu-mono/internal/config"
	"github.com/yuldoshevgg/menu-mono/internal/core/repository"
	"github.com/yuldoshevgg/menu-mono/internal/pkg/logger"
	"github.com/yuldoshevgg/menu-mono/internal/transport/grpc"
	handlers "github.com/yuldoshevgg/menu-mono/internal/transport/handler"
	grpcMain "google.golang.org/grpc"
)

var (
	cfg *config.Config
)

func init() {
	log.Println("Initializing...")
	cfg = config.Load()
	logger.SetLogger(&cfg.Logging)
	logger.Log.Info("Initializing done...")
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Project.Timeout)
	defer cancel()

	repos := repository.New(ctx, cfg.PSQL.URI)

	grpcServer := grpc.New(&grpc.GrpcTransport{
		Cfg:  cfg,
		Repo: repos,
	})

	lis, err := net.Listen("tcp", cfg.Grpc.URL)
	if err != nil {
		panic(err)
	}

	go func() {
		logger.Log.Info("starting grpc server on " + cfg.Grpc.URL)
		if err = grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	router := setUpHTTP(repos)
	go func() {
		logger.Log.Info("starting http server on " + cfg.HTTP.URL)
		if err := router.Run(cfg.HTTP.URL); err != nil {
			panic(err)
		}
	}()

	gracefulShutdown(grpcServer, ctx, cancel)
}

func gracefulShutdown(grpcServer *grpcMain.Server, ctx context.Context, cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		logger.Log.Info("shutting down")

		grpcServer.GracefulStop()

		logger.Log.Info("shutdown successfully called")
		wg.Done()
	}(&wg)

	go func() {
		wg.Wait()
		cancel()
	}()

	<-ctx.Done()
	os.Exit(0)
}

func setUpHTTP(repos repository.Store) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	switch cfg.Mode {
	case "dev":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "*")

	router.Use(cors.New(corsConfig))
	v1 := router.Group("/v1")

	gwMux := handlers.New(context.Background(), cfg)
	router.Static("/swagger", "./doc/swagger")
	v1.Any("/*any", func(c *gin.Context) {
		gwMux.ServeHTTP(c.Writer, c.Request)
	})

	return router
}
