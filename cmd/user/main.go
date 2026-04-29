package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	authv1 "github.com/qkitzero/auth-service/gen/go/auth/v1"
	userv1 "github.com/qkitzero/user-service/gen/go/user/v1"
	appuser "github.com/qkitzero/user-service/internal/application/user"
	apiauth "github.com/qkitzero/user-service/internal/infrastructure/api/auth"
	"github.com/qkitzero/user-service/internal/infrastructure/db"
	infrauser "github.com/qkitzero/user-service/internal/infrastructure/user"
	grpcuser "github.com/qkitzero/user-service/internal/interface/grpc/user"
)

const shutdownTimeout = 15 * time.Second

type config struct {
	Env             string
	Port            string
	DBHost          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBPort          string
	DBSSLMode       string
	AuthServiceHost string
	AuthServicePort string
}

func loadConfig() (config, error) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	cfg := config{Env: env}
	required := []struct {
		key string
		dst *string
	}{
		{"PORT", &cfg.Port},
		{"DB_HOST", &cfg.DBHost},
		{"DB_USER", &cfg.DBUser},
		{"DB_PASSWORD", &cfg.DBPassword},
		{"DB_NAME", &cfg.DBName},
		{"DB_PORT", &cfg.DBPort},
		{"DB_SSL_MODE", &cfg.DBSSLMode},
		{"AUTH_SERVICE_HOST", &cfg.AuthServiceHost},
		{"AUTH_SERVICE_PORT", &cfg.AuthServicePort},
	}
	var missing []string
	for _, r := range required {
		v := os.Getenv(r.key)
		if v == "" {
			missing = append(missing, r.key)
			continue
		}
		*r.dst = v
	}
	if len(missing) > 0 {
		return cfg, fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}
	return cfg, nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("user-service: %v", err)
	}
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	gormDB, err := db.Init(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)
	if err != nil {
		return fmt.Errorf("db init: %w", err)
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("db handle: %w", err)
	}
	defer func() { _ = sqlDB.Close() }()

	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	var dialOpt grpc.DialOption
	switch cfg.Env {
	case "production":
		dialOpt = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	default:
		dialOpt = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	conn, err := grpc.NewClient(cfg.AuthServiceHost+":"+cfg.AuthServicePort, dialOpt)
	if err != nil {
		return fmt.Errorf("auth client: %w", err)
	}
	defer func() { _ = conn.Close() }()

	server := grpc.NewServer()

	authServiceClient := authv1.NewAuthServiceClient(conn)
	userRepository := infrauser.NewUserRepository(gormDB)

	authService := apiauth.NewAuthService(authServiceClient)
	userUsecase := appuser.NewUserUsecase(authService, userRepository)

	healthServer := health.NewServer()
	userHandler := grpcuser.NewUserHandler(userUsecase)

	grpc_health_v1.RegisterHealthServer(server, healthServer)
	userv1.RegisterUserServiceServer(server, userHandler)

	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("user", grpc_health_v1.HealthCheckResponse_SERVING)

	if cfg.Env == "development" {
		reflection.Register(server)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serveErr := make(chan error, 1)
	go func() {
		log.Printf("gRPC server listening on %s", listener.Addr().String())
		serveErr <- server.Serve(listener)
	}()

	select {
	case err := <-serveErr:
		if err != nil {
			return fmt.Errorf("grpc serve: %w", err)
		}
		return nil
	case <-ctx.Done():
		log.Println("shutdown signal received, starting graceful stop")
		healthServer.Shutdown()

		stopped := make(chan struct{})
		go func() {
			server.GracefulStop()
			close(stopped)
		}()

		select {
		case <-stopped:
			log.Println("gRPC server stopped gracefully")
		case <-time.After(shutdownTimeout):
			log.Printf("graceful stop timed out after %s, forcing stop", shutdownTimeout)
			server.Stop()
			<-stopped
		}
		return nil
	}
}
