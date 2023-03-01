package client

import (
	"context"
	"go-usermgmt-grpc/pb"
	"time"

	"google.golang.org/grpc"
)

type AuthClient struct {
	service  pb.AccountsClient
	username string
	password string
}

func NewAuthClient(cc *grpc.ClientConn, username string, password string) *AuthClient {
	serive := pb.NewAccountsClient(cc)
	return &AuthClient{
		service:  serive,
		username: username,
		password: password,
	}
}

func (client *AuthClient) Login() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.LoginRequest{
		Username: client.username,
		Password: client.password,
	}
	res, err := client.service.LoginUser(ctx, req)
	if err != nil {
		return "", err
	}
	return res.GetAccessToken(), nil
}
