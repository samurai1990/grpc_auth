package service

import (
	"context"
	"go-usermgmt-grpc/db"
	"go-usermgmt-grpc/pb"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	SecretKey        = "secrettest"
	TokenDuration    = 15 * time.Minute
	minSecretKeySize = 32
)

type UserManagementServer struct {
	pb.UnimplementedAccountsServer
}

func (ums *UserManagementServer) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Printf("Received: %v", in.GetUser()) // test
	storeUser := &db.UserStore{}

	Newuser, err := CreateUser(
		storeUser,
		in.User.GetUsername(),
		in.User.GetPassword(),
		in.User.GetEmail(),
		in.User.GetIsAdmin(),
		in.User.GetIsActive(),
	)
	if err != nil && err == db.ErrAlreadyExists {
		return nil, status.Errorf(codes.AlreadyExists, "user: '%s' is exists", in.User.GetUsername())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot save user: %v", err)
	}
	userIDToString := Newuser.Id.String()
	return &pb.CreateUserResponse{
		Id:       userIDToString,
		Username: Newuser.Username,
		Email:    Newuser.Email,
		IsAdmin:  Newuser.IsAdmin,
		IsActive: Newuser.IsActive,
	}, nil
}

func (ums *UserManagementServer) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	storeUser := &db.UserStore{}
	user, err := storeUser.Find(req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}

	if user == nil || !user.IsCorrectPassword(req.GetPassword()) {
		return nil, status.Errorf(codes.NotFound, "incorect username/password")
	}
	jwtManager := JWTManager{
		secretKey:     SecretKey,
		tokenDuration: TokenDuration,
	}
	token, err := jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cnanot generate access token")
	}

	res := &pb.LoginResponse{AccessToken: token}
	return res, nil
}

func (ums *UserManagementServer) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserRespose, error) {
	userStore := db.UserStore{}
	objUsers, err := userStore.List()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "no list usrs")
	}

	listUser := []*pb.User{}
	for _, user := range objUsers {
		listUser = append(listUser, &pb.User{
			Username: user.Username,
			Email:    user.Email,
			IsAdmin:  user.IsAdmin,
			IsActive: user.IsActive,
			Id:       user.Id.String(),
		})
	}
	return &pb.ListUserRespose{
		Users: listUser,
	}, nil
}
