package main

import (
	"log"
	"net"

	auth_pb "github.com/qkitzero/auth/pb"
	userv1 "github.com/qkitzero/user/gen/go/proto/user/v1"
	application_user "github.com/qkitzero/user/internal/application/user"
	api_auth "github.com/qkitzero/user/internal/infrastructure/api/auth"
	"github.com/qkitzero/user/internal/infrastructure/db"
	infrastructure_user "github.com/qkitzero/user/internal/infrastructure/user"
	interface_user "github.com/qkitzero/user/internal/interface/grpc/user"
	"github.com/qkitzero/user/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	db, err := db.Init(
		util.GetEnv("DB_USER"),
		util.GetEnv("DB_PASSWORD"),
		util.GetEnv("DB_HOST"),
		util.GetEnv("DB_PORT"),
		util.GetEnv("DB_NAME"),
	)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", ":"+util.GetEnv("PORT"))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.NewClient(
		util.GetEnv("AUTH_SERVICE_HOST")+":"+util.GetEnv("AUTH_SERVICE_PORT"),
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

	userv1.RegisterUserServiceServer(server, userHandler)

	if err = server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
