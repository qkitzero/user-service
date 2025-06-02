package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"

	userv1 "github.com/qkitzero/user/gen/go/user/v1"
	"github.com/qkitzero/user/util"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		util.GetEnv("SERVER_HOST", "")+":"+util.GetEnv("SERVER_PORT", ""),
		grpc.WithTransportCredentials(insecure.NewCredentials()), // dev
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	healthClient := grpc_health_v1.NewHealthClient(conn)

	mux := runtime.NewServeMux(
		runtime.WithHealthzEndpoint(healthClient),
	)
	endpoint := util.GetEnv("SERVER_HOST", "") + ":" + util.GetEnv("SERVER_PORT", "")
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := userv1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":"+util.GetEnv("GRPC_GATEWAY_PORT", ""), mux); err != nil {
		log.Fatal(err)
	}
}
