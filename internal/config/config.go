package config

import (
	"fmt"
	"log"
	"os"
	"time"

	logging "github.com/yuldoshevgg/menu-mono/internal/pkg/logger"
	"github.com/yuldoshevgg/menu-mono/pkg/logger"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type AppMode string

const (
	DEVELOPMENT AppMode = "DEVELOPMENT"
	PRODUCTION  AppMode = "PRODUCTION"

	GrpcGatewayHTTPPath   = "grpc-gateway-http-path"
	GrpcGatewayHTTPMethod = "grpc-gateway-http-method"
	GrpcReportID          = "X-Report-Id"
	OneDay                = time.Hour * 24
)

var (
	TimeoutDuration time.Duration
	CacheTimeout    time.Duration
)

type Config struct {
	Logging logger.LoggingConfig
	Mode    string `env:"APPLICATION_MODE"`

	Project struct {
		Name                string        `env:"PROJECT_NAME"`
		Version             string        `env:"PROJECT_VERSION"`
		Timeout             time.Duration `env:"PROJECT_TIMEOUT"`
		ProviderInitTimeout time.Duration `env:"PROJECT_TIMEOUT"`
		SwaggerEnabled      bool          `env:"PROJECT_SWAGGER_ENABLED"`
		CacheTimeout        time.Duration `env:"PROJECT_CACHE_TIMEOUT"`
	}

	App struct {
		BaseURL string `env:"APP_BASE_URL"`
	}

	Token struct {
		AccessExpiresInTime  time.Duration `env:"TOKEN_ACCESS_EXPIRES_IN_TIME"`
		TokenExpiresInTime   time.Duration `env:"TOKEN_EXPIRES_IN_TIME"`
		RefreshExpiresInTime time.Duration `env:"TOKEN_REFRESH_EXPIRES_IN_TIME"`
		SecretKey            string        `env:"TOKEN_SECRET_KEY"`
	}

	HTTP struct {
		Host string `env:"HTTP_HOST"`
		Port int    `env:"HTTP_PORT"`

		URL string `env:"HTTP_URL"`
	}

	Grpc struct {
		Host string `env:"GRPC_HOST"`
		Port int    `env:"GRPC_PORT"`

		URL string `env:"GRPC_URL"`
	}

	PSQL struct {
		URI string `env:"PSQL_URI"`
	}
}

func Load() *Config {
	var (
		cfg  Config
		path = ".env"
	)

	if err := godotenv.Load(path); err != nil && !os.IsNotExist(err) {
		logging.Log.Warn("failed loading .env file", zap.Error(err))
	}

	if err := env.Parse(&cfg); err != nil {
		log.Println(err.Error())
		panic("unmarshal from environment error")
	}

	TimeoutDuration = cfg.Project.Timeout

	cfg.MakeGrpcURL()
	cfg.MakeHTTPURL()

	return &cfg
}

func (c *Config) MakeGrpcURL() {
	c.Grpc.URL = fmt.Sprintf("%s:%d", c.Grpc.Host, c.Grpc.Port)
}

func (c *Config) MakeHTTPURL() {
	c.HTTP.URL = fmt.Sprintf("%s:%d", c.HTTP.Host, c.HTTP.Port)
}
