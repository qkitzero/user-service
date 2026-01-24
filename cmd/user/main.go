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
	"gorm.io/gorm"

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := db.Init(
		util.GetEnv("DB_HOST", "localhost"),
		util.GetEnv("DB_USER", "user"),
		util.GetEnv("DB_PASSWORD", "password"),
		util.GetEnv("DB_NAME", "user_db"),
		util.GetEnv("DB_PORT", "5432"),
		util.GetEnv("DB_SSL_MODE", "disable"),
	)
	if err != nil {
		log.Fatal(err)
	}

	userHandler, authConn := newDependencies(db)
	defer authConn.Close()

	grpcPort := util.GetEnv("GRPC_PORT", "50051")
	grpcServer := newGRPCServer(userHandler)

	go startGRPCServer(grpcServer, grpcPort)

	httpPort := util.GetEnv("HTTP_PORT", "8080")
	startHTTPServer(ctx, grpcPort, httpPort)
}

func newDependencies(db *gorm.DB) (*grpcuser.UserHandler, *grpc.ClientConn) {
	authTarget := util.GetEnv("AUTH_SERVICE_HOST", "localhost") + ":" + util.GetEnv("AUTH_SERVICE_PORT", "50051")

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

	authServiceClient := authv1.NewAuthServiceClient(authConn)
	userRepository := infrauser.NewUserRepository(db)

	authUsecase := apiauth.NewAuthUsecase(authServiceClient)
	userUsecase := appuser.NewUserUsecase(userRepository)

	userHandler := grpcuser.NewUserHandler(authUsecase, userUsecase)

	return userHandler, authConn
}

func newGRPCServer(userHandler *grpcuser.UserHandler) *grpc.Server {
	healthServer := health.NewServer()
	healthServer.SetServingStatus("user", grpc_health_v1.HealthCheckResponse_SERVING)

	server := grpc.NewServer()

	grpc_health_v1.RegisterHealthServer(server, healthServer)
	userv1.RegisterUserServiceServer(server, userHandler)

	if util.GetEnv("ENV", "development") == "development" {
		reflection.Register(server)
	}

	return server
}

func startGRPCServer(server *grpc.Server, port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("starting gRPC server on port %s", port)

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func startHTTPServer(ctx context.Context, grpcPort, httpPort string) {
	grpcEndpoint := "localhost:" + grpcPort

	userConn, err := grpc.NewClient(grpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer userConn.Close()

	mux := runtime.NewServeMux(
		runtime.WithHealthzEndpoint(grpc_health_v1.NewHealthClient(userConn)),
	)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := userv1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts); err != nil {
		log.Fatal(err)
	}

	log.Printf("starting http server on port %s", httpPort)

	if err := http.ListenAndServe(":"+httpPort, mux); err != nil {
		log.Fatal(err)
	}
}
