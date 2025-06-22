package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	authv1 "github.com/qkitzero/auth/gen/go/auth/v1"
	userv1 "github.com/qkitzero/user/gen/go/user/v1"
	application_user "github.com/qkitzero/user/internal/application/user"
	api_auth "github.com/qkitzero/user/internal/infrastructure/api/auth"
	"github.com/qkitzero/user/internal/infrastructure/db"
	infrastructure_user "github.com/qkitzero/user/internal/infrastructure/user"
	interface_user "github.com/qkitzero/user/internal/interface/grpc/user"
	"github.com/qkitzero/user/util"
)

func main() {
	db, err := db.Init(
		util.GetEnv("DB_HOST", ""),
		util.GetEnv("DB_USER", ""),
		util.GetEnv("DB_PASSWORD", ""),
		util.GetEnv("DB_NAME", ""),
		util.GetEnv("DB_PORT", ""),
		util.GetEnv("DB_SSL_MODE", ""),
	)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", ":"+util.GetEnv("PORT", ""))
	if err != nil {
		log.Fatal(err)
	}

	authTarget := util.GetEnv("AUTH_SERVICE_HOST", "") + ":" + util.GetEnv("AUTH_SERVICE_PORT", "")

	var opts grpc.DialOption
	switch util.GetEnv("ENV", "development") {
	case "production":
		opts = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	default:
		opts = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	conn, err := grpc.NewClient(authTarget, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	server := grpc.NewServer()

	authServiceClient := authv1.NewAuthServiceClient(conn)
	userRepository := infrastructure_user.NewUserRepository(db)

	authUsecase := api_auth.NewAuthUsecase(authServiceClient)
	userUsecase := application_user.NewUserUsecase(userRepository)

	healthServer := health.NewServer()
	userHandler := interface_user.NewUserHandler(authUsecase, userUsecase)

	grpc_health_v1.RegisterHealthServer(server, healthServer)
	userv1.RegisterUserServiceServer(server, userHandler)

	healthServer.SetServingStatus("user", grpc_health_v1.HealthCheckResponse_SERVING)

	if util.GetEnv("ENV", "development") == "development" {
		reflection.Register(server)
	}

	if err = server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
