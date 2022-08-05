package clientservices

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/clientservices/httpclient"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/gtools"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userpb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userpb/userpbmapper"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserService interface {
	GetByAuthId(ctx context.Context, authId string) single.Single[userdtos.UserDto]
	GetByAuthIdAsync(ctx context.Context, authId string) single.Single[userdtos.UserDto]
}

type UserServiceImpl struct {
	systemAccessTokenClient httpclient.SysAccessTokenClient
}

func (u UserServiceImpl) GetByAuthIdAsync(ctx context.Context, authId string) single.Single[userdtos.UserDto] {
	return single.MapIdentityAsync[userdtos.UserDto](ctx, u.GetByAuthId(ctx, authId))
}
func (u UserServiceImpl) GetByAuthId(ctx context.Context, authId string) single.Single[userdtos.UserDto] {
	accessTokenSrc := single.FromSupplier(u.systemAccessTokenClient.GetGRPCAccessToken)
	connectionSrc := single.MapWithError(accessTokenSrc, func(token oauth2.Token) (*grpc.ClientConn, error) {
		//perRPC := oauth.NewOauthAccess(&token)
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			//grpc.WithPerRPCCredentials(perRPC),
		}
		return grpc.Dial("localhost:50051", opts...)
	})
	userReplySrc := single.MapWithError(connectionSrc, func(conn *grpc.ClientConn) (*userpb.UserReply, error) {
		defer conn.Close()
		userService := userpb.NewUserServiceClient(conn)
		reply, err := userService.GetUserByAuthId(ctx, &userpb.AuthIdRequest{AuthId: authId})
		return reply, gtools.NewErrorResponseHandler(err).GetProcessedError()
	})
	return single.Map(userReplySrc, func(reply *userpb.UserReply) userdtos.UserDto {
		dto := userdtos.UserDto{}
		userpbmapper.MapUserReplyToUserDto(reply, &dto)
		return dto
	})
}

func NewUserService(systemAccessTokenClient httpclient.SysAccessTokenClient) *UserServiceImpl {

	return &UserServiceImpl{systemAccessTokenClient: systemAccessTokenClient}
}
