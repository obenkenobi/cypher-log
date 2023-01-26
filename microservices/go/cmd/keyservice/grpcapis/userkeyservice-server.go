package grpcapis

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/gtools"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userkeypb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/commondtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmappers/grpcmappers"
)

type UserKeyServiceServerImpl struct {
	userkeypb.UnimplementedUserKeyServiceServer
	userKeyService services.UserKeyService
}

func (u UserKeyServiceServerImpl) GetKeyFromSession(
	ctx context.Context,
	userKeySession *userkeypb.UserKeySession,
) (*userkeypb.UserKey, error) {
	userKeySessionDto := commondtos.UKeySessionDto{}
	grpcmappers.UserKeySessionToUserKeySessionDto(userKeySession, &userKeySessionDto)
	keyDto, err := u.userKeyService.GetKeyFromSession(ctx, userKeySessionDto)
	if err != nil {
		return nil, gtools.ProcessErrorToGrpcStatusError(ctx, gtools.ReadAction, err)
	}
	userKey := &userkeypb.UserKey{}
	grpcmappers.UserKeyDtoToUserKey(&keyDto, userKey)
	return userKey, nil
}

func NewUserKeyServiceServerImpl(userKeyService services.UserKeyService) *UserKeyServiceServerImpl {
	return &UserKeyServiceServerImpl{userKeyService: userKeyService}
}
