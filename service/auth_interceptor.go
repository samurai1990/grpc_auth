package service

import (
	"context"
	"go-usermgmt-grpc/db"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwt            *JWTManager
	accessibleRole map[string][]string
}

func NewAuthInterceptor(jwtManager *JWTManager, accessibleRole map[string][]string) *AuthInterceptor {
	return &AuthInterceptor{
		jwt:            jwtManager,
		accessibleRole: accessibleRole,
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) error {
	accessibleRoles, ok := interceptor.accessibleRole[method]
	if !ok {
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	claim, err := interceptor.jwt.Verify(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}
	userStore := &db.UserStore{}
	queryUser := FindUser(userStore, claim.Username)
	if queryUser == nil {
		return status.Errorf(codes.NotFound, "user not found")
	}
	if queryUser.IsAdmin != claim.IsAdmin {
		return status.Errorf(codes.InvalidArgument, "malformed token")
	}
	for _, role := range accessibleRoles {
		if role == "is_admin" {
			if queryUser.IsAdmin {
				return nil
			}
		}
	}
	return status.Error(codes.PermissionDenied, "no permission to access this RPC")
}

func (interceptor *AuthInterceptor) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("--> unary interceptor: ", info.FullMethod)

	err := interceptor.authorize(ctx, info.FullMethod)
	if err != nil {
		return nil, err
	}
	return handler(ctx, req)
}
