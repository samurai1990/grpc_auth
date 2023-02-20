package main

import (
	"flag"
	"go-usermgmt-grpc/client"

	"go-usermgmt-grpc/pb"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	username        = "admin1"
	password        = "admin"
	refreshDuration = 30 * time.Second
)

func testCreateuser(userClient *client.UserClient) {
	newUser := &pb.Newuser{
		Name: "piter",
		Age:  25,
	}
	userClient.CreateUser(newUser)
}

func authMethods() map[string]bool {
	const userManagePath = "/usermgmt.UserManagement/"
	return map[string]bool{
		userManagePath + "CreateNewUser": true,
		userManagePath + "ListNewUser":   true,
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

	authClient := client.NewAuthClient(conn, username, password)
	interceptor, err := client.NewAuthInterceptor(authClient, authMethods(), refreshDuration)
	if err != nil {
		log.Fatal("connot create auth interceptor: ", err)
	}
	cc2, err := grpc.Dial(
		*serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
	)
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}
	userClient := client.NewUserClient(cc2)
	testCreateuser(userClient)
}
