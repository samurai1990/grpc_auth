package client

import (
	"context"
	"go-usermgmt-grpc/pb"
	"time"

	"google.golang.org/grpc"
)

type UserClient struct {
	userClient pb.AccountsClient
}

func NewUserClient(cc *grpc.ClientConn) *UserClient {
	userClient := pb.NewAccountsClient(cc)
	return &UserClient{userClient: userClient}
}

func (userClient *UserClient) ListUser() (*pb.ListUserRespose, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &pb.ListUserRequest{}
	res, err := userClient.userClient.ListUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
