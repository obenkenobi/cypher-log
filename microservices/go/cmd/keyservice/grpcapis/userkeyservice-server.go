package grpcapis

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/gtools"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userkeypb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmappers/grpcmappers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/keydtos"
)

type UserKeyServiceServerImpl struct {
	userkeypb.UnimplementedUserKeyServiceServer
	userKeyService services.UserKeyService
}

func (u UserKeyServiceServerImpl) GetKeyFromSession(
	ctx context.Context,
	userKeySession *userkeypb.UserKeySession,
) (*userkeypb.UserKey, error) {
	userKeySessionDto := keydtos.UserKeySessionDto{}
	grpcmappers.UserKeySessionToUserKeySessionDto(userKeySession, &userKeySessionDto)
	keySrc := u.userKeyService.GetKeyFromSession(ctx, userKeySessionDto)
	replySrc := single.Map(keySrc, func(keyDto keydtos.UserKeyDto) *userkeypb.UserKey {
		userKey := &userkeypb.UserKey{}
		grpcmappers.UserKeyDtoToUserKey(&keyDto, userKey)
		return userKey
	})
	res, err := single.RetrieveValue(ctx, replySrc)
	return res, gtools.ProcessErrorToGrpcStatusError(gtools.ReadAction, err)
}

func NewUserKeyServiceServerImpl(userKeyService services.UserKeyService) *UserKeyServiceServerImpl {
	return &UserKeyServiceServerImpl{userKeyService: userKeyService}
}
