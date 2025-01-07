package main

import (
	"fmt"
	"log"
	"net"
	"os"

	application_user "user/internal/application/user"
	"user/internal/infrastructure/db"
	infrastructure_user "user/internal/infrastructure/persistence/user"
	interface_user "user/internal/interface/grpc/user"
	"user/pb"

	"google.golang.org/grpc"
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

	server := grpc.NewServer()

	userRepository := infrastructure_user.NewUserRepository(db)

	userService := application_user.NewUserService(userRepository)

	userHandler := interface_user.NewUserHandler(userService)

	pb.RegisterUserServiceServer(server, userHandler)

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
