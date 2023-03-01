package main

import (
	"flag"
	"fmt"
	"go-usermgmt-grpc/pb"
	"go-usermgmt-grpc/service"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func accessibleRoles() map[string][]string {
	const userManagePath = "/pb.Accounts/"

	return map[string][]string{

		userManagePath + "CreateUser": {"is_admin"},
		userManagePath + "ListUser":     {"is_admin"},
	}
}

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("start server on port %d:", *port)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	usermanager := service.UserManagementServer{}
	jwtManager := service.NewJWTManager(service.SecretKey, service.TokenDuration)
	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.UnaryInterceptor),
	)
	reflection.Register(grpcServer) // evans
	pb.RegisterAccountsServer(grpcServer, &usermanager)
	log.Printf("server listening at %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC server over %s: %v", address, err)
	}
}
