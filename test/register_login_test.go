package tests

import (
	"context"
	"flag"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/serzap/auth_service/api"
	"github.com/serzap/auth_service/internal/config"
	"github.com/serzap/auth_service/internal/server"
	"github.com/serzap/auth_service/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "../etc/test.yaml", "the config file")

func TestGrpc(t *testing.T) {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		api.RegisterAuthServiceServer(grpcServer, server.NewAuthServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	go func() {
		fmt.Printf("Starting gRPC server at %s...\n", c.ListenOn)
		s.Start()
	}()
	defer s.Stop()

	time.Sleep(time.Second)

	clientCtx := context.Background()
	log.Println(c.ListenOn)
	clientConf := zrpc.RpcClientConf{
		Target: c.ListenOn,
	}
	client := api.NewAuthServiceClient(zrpc.MustNewClient(clientConf).Conn())

	registerRequest := &api.RegisterRequest{

		Email:     "test@example.com",
		Password:  "password123",
		Username:  "test",
		FirstName: "test",
		LastName:  "test",
	}

	registerResponse, err := client.Register(clientCtx, registerRequest)
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}
	fmt.Printf("Register Response: %v\n", registerResponse)

	loginRequest := &api.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	loginResponse, err := client.Login(clientCtx, loginRequest)
	if err != nil {
		log.Fatalf("could not login: %v", err)
	}
	fmt.Printf("Login Response: %v\n", loginResponse)
}
