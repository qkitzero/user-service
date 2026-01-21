package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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
	"github.com/qkitzero/user-service/util"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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

	authTarget := util.GetEnv("AUTH_SERVICE_HOST", "") + ":" + util.GetEnv("AUTH_SERVICE_PORT", "")

	var opts grpc.DialOption
	switch util.GetEnv("ENV", "development") {
	case "production":
		opts = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	default:
		opts = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	authConn, err := grpc.NewClient(authTarget, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer authConn.Close()

	authServiceClient := authv1.NewAuthServiceClient(authConn)
	userRepository := infrauser.NewUserRepository(db)

	authUsecase := apiauth.NewAuthUsecase(authServiceClient)
	userUsecase := appuser.NewUserUsecase(userRepository)

	userHandler := grpcuser.NewUserHandler(authUsecase, userUsecase)

	healthServer := health.NewServer()
	healthServer.SetServingStatus("user", grpc_health_v1.HealthCheckResponse_SERVING)

	grpcServer := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	userv1.RegisterUserServiceServer(grpcServer, userHandler)

	if util.GetEnv("ENV", "development") == "development" {
		reflection.Register(grpcServer)
	}

	grpcPort := util.GetEnv("GRPC_PORT", "50051")
	grpcListener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Printf("starting grpc server on port %s", grpcPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatal(err)
		}
	}()

	userConn, err := grpc.NewClient("localhost:"+grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer userConn.Close()

	mux := runtime.NewServeMux(
		runtime.WithHealthzEndpoint(grpc_health_v1.NewHealthClient(userConn)),
	)

	if err := userv1.RegisterUserServiceHandlerServer(ctx, mux, userHandler); err != nil {
		log.Fatal(err)
	}

	httpPort := util.GetEnv("HTTP_PORT", "8080")
	log.Printf("starting http server on port %s", httpPort)
	if err := http.ListenAndServe(":"+httpPort, mux); err != nil {
		log.Fatal(err)
	}
}
