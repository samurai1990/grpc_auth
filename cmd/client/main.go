package main

import (
	"context"
	"flag"
	"go-usermgmt-grpc/pb"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatal("cannnot dial server: ", err)
	}
	defer conn.Close()

	c := pb.NewUserManagementClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var new_users = make(map[string]int32)
	new_users["alice"] = 43
	new_users["bob"] = 30
	for name, age := range new_users {
		r, err := c.CreateNewUser(ctx, &pb.Newuser{Name: name, Age: age})
		if err != nil {
			log.Fatal("could not create user: ", err)
		}
		log.Printf(`User Details:
		NAME: %s
		AGE: %d
		ID: %d`, r.GetName(), r.GetAge(), r.GetId())
	}

}
