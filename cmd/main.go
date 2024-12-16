package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	pb "user/pb"

	"google.golang.org/grpc"
)

type myServer struct {
	pb.UnimplementedUserServiceServer
}

func NewMyServer() *myServer {
	return &myServer{}
}

func main() {
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	pb.RegisterUserServiceServer(s, NewMyServer())

	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}

func (s *myServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return &pb.CreateUserResponse{
		Message: fmt.Sprintf("Create User, %s!", req.GetName()),
	}, nil
}
