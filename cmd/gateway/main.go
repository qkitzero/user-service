package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userv1 "github.com/qkitzero/user/gen/go/proto/user/v1"
	"github.com/qkitzero/user/util"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	endpoint := util.GetEnv("SERVER_HOST") + ":" + util.GetEnv("SERVER_PORT")
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := userv1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":"+util.GetEnv("GRPC_GATEWAY_PORT"), mux); err != nil {
		log.Fatal(err)
	}
}
