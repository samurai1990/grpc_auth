package main

import (
	"context"
	"flag"
	"fmt"
	"go-usermgmt-grpc/pb"
	"go-usermgmt-grpc/service"
	"log"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
}

func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.Newuser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())
	var user_id int32 = int32(rand.Intn(1000))
	return &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}, nil
}

func unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("--> unary interceptor: ", info.FullMethod)
	return handler(ctx, req)
}

const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
)

func seedUser(userStore service.UserStore) error {
	err := createUser(userStore, "admin1", "admin", "admin")
	if err != nil {
		return err
	}
	return createUser(userStore, "user1", "user", "user")

}

func createUser(userStore service.UserStore, username string, password string, role string) error {
	user, err := service.NewUser(username, password, role)
	if err != nil {
		return nil
	}
	return userStore.Save(user)
}

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("start server on port %d:", *port)

	userStore := service.NewInMemoryUserStore()
	jwtManager := service.NewJWTManager(secretKey, tokenDuration)
	authServer := service.NewAuthServer(userStore, jwtManager)

	err := seedUser(userStore)
	if err != nil {
		log.Fatal("cannot seed users")
	}

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor),
	)
	reflection.Register(grpcServer)
	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterUserManagementServer(grpcServer, &UserManagementServer{})
	log.Printf("server listening at %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
