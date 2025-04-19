package main

import (
	"fmt"
	"log"
	"net"
	"os"

	auth_pb "github.com/qkitzero/auth/pb"
	application_user "github.com/qkitzero/user/internal/application/user"
	api_auth "github.com/qkitzero/user/internal/infrastructure/api/auth"
	"github.com/qkitzero/user/internal/infrastructure/db"
	infrastructure_user "github.com/qkitzero/user/internal/infrastructure/user"
	interface_user "github.com/qkitzero/user/internal/interface/grpc/user"
	user_pb "github.com/qkitzero/user/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	db, err := db.Init(
		getEnv("DB_USER"),
		getEnv("DB_PASSWORD"),
		getEnv("DB_HOST"),
		getEnv("DB_PORT"),
		getEnv("DB_NAME"),
	)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", ":"+getEnv("PORT"))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.NewClient(
		getEnv("AUTH_SERVICE_HOST")+":"+getEnv("AUTH_SERVICE_PORT"),
		grpc.WithTransportCredentials(insecure.NewCredentials()), // dev
	)
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()

	authServiceClient := auth_pb.NewAuthServiceClient(conn)
	userRepository := infrastructure_user.NewUserRepository(db)

	authUsecase := api_auth.NewAuthUsecase(authServiceClient)
	userUsecase := application_user.NewUserUsecase(userRepository)

	userHandler := interface_user.NewUserHandler(authUsecase, userUsecase)

	user_pb.RegisterUserServiceServer(server, userHandler)

	if err = server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatal(fmt.Sprintf("missing required environment variable: %s", key))
	}
	return value
}
