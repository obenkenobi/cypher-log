package grpcserveropts

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/securityservices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

type AuthInterceptorCreator interface {
	CreateUnaryInterceptor() grpc.ServerOption
}

type AuthInterceptorCreatorImpl struct {
	grpcAuth0JwtValidateService securityservices.ExternalOath2ValidateService
}

// valid validates the authorization.
func (a AuthInterceptorCreatorImpl) valid(ctx context.Context, authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	_, err := a.grpcAuth0JwtValidateService.GetJwtValidator().ValidateToken(ctx, token)
	return err == nil
}

func (a AuthInterceptorCreatorImpl) authenticate(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	if !a.valid(ctx, md["authorization"]) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}

func (a AuthInterceptorCreatorImpl) CreateUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(a.authenticate)
}

func NewAuthInterceptorCreator(
	grpcAuth0JwtValidateService securityservices.ExternalOath2ValidateService,
) AuthInterceptorCreator {
	return &AuthInterceptorCreatorImpl{grpcAuth0JwtValidateService: grpcAuth0JwtValidateService}
}
