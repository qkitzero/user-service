package main

import (
	"fmt"
	"log"
	"net"
	"os"

	application_user "user/internal/application/user"
	"user/internal/infrastructure/api"
	"user/internal/infrastructure/db"
	infrastructure_user "user/internal/infrastructure/persistence/user"
	interface_user "user/internal/interface/grpc/user"
	user_pb "user/pb"

	auth_pb "github.com/qkitzero/auth/pb"
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

	authService := api.NewAuthService(authServiceClient)
	userService := application_user.NewUserService(userRepository)

	userHandler := interface_user.NewUserHandler(authService, userService)

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
