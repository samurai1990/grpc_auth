package client

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	authMethod  map[string]bool
	accessToken string
}

func NewAuthInterceptor(authMethod map[string]bool, accessToken string) *AuthInterceptor {
	return &AuthInterceptor{
		authMethod:  authMethod,
		accessToken: accessToken,
	}

}

func (interceptor *AuthInterceptor) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {
		log.Printf("--> unary interceptor: %s", method)
		if interceptor.authMethod[method] {
			return invoker(interceptor.attachToken(ctx), method, req, reply, cc, opts...)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (interceptor *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", interceptor.accessToken)
}
