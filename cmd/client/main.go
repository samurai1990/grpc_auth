package main

import (
	"flag"
	"go-usermgmt-grpc/client"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	username        = "test"
	password        = "test"
	refreshDuration = 30 * time.Second
)

func authMethods() map[string]bool {
	const userManagePath = "/pb.Accounts/"
	return map[string]bool{
		userManagePath + "CreateNewUser": true,
		userManagePath + "ListUser":      true,
	}
}

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatal("cannnot dial server: ", err)
	}
	defer conn.Close()

	authClient := client.NewAuthClient(conn, username, password)
	token, err := authClient.Login()
	if err != nil {
		log.Fatal("connot create auth interceptor: ", err)
	}
	log.Printf("accessToken:%s", token)

	authInterceptor := client.NewAuthInterceptor(authMethods(), token)
	cc2, err := grpc.Dial(*serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(authInterceptor.UnaryInterceptor()),
	)
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}
	defer cc2.Close()

	userClient := client.NewUserClient(cc2)
	users, err := userClient.ListUser()
	if err != nil {
		log.Fatal("cannot list users with error: ", err)
	}
	log.Print(users)

}
