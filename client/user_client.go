package client

import (
	"context"
	"go-usermgmt-grpc/pb"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserClient struct {
	service pb.UserManagementClient
}

func NewUserClient(cc *grpc.ClientConn) *UserClient {
	service := pb.NewUserManagementClient(cc)
	return &UserClient{service: service}
}

func (userClient *UserClient) CreateUser(user *pb.Newuser) {
	req := &pb.Newuser{
		Name: user.Name,
		Age:  user.Age,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := userClient.service.CreateNewUser(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Print("user already exsist")
		} else {
			log.Fatal("cannot create user: ", err)
		}
		return
	}
	log.Printf("created user id: %d", res.Id)
	log.Printf(`User Details:
	NAME: %s
	AGE: %d
	ID: %d`, res.GetName(), res.GetAge(), res.GetId())
}
