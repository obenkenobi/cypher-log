package clientservices

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
	env "github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/gtools"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userpb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userpb/userpbmapper"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"google.golang.org/grpc"
)

type UserService interface {
	GetByAuthId(ctx context.Context, authId string) single.Single[userdtos.UserDto]
	GetByAuthIdAsync(ctx context.Context, authId string) single.Single[userdtos.UserDto]
}

type UserServiceImpl struct {
	systemAccessTokenClient SysAccessTokenClient
	tlsConf                 conf.TLSConf
}

func (u UserServiceImpl) GetByAuthIdAsync(ctx context.Context, authId string) single.Single[userdtos.UserDto] {
	return single.MapIdentityAsync[userdtos.UserDto](ctx, u.GetByAuthId(ctx, authId))
}
func (u UserServiceImpl) GetByAuthId(ctx context.Context, authId string) single.Single[userdtos.UserDto] {
	optsSrc := gtools.CreateSingleWithDialOptionsIfAuthActivated(
		env.ActivateGRPCAuth(),
		[]gtools.DialOptionSingleCreator{
			func() single.Single[grpc.DialOption] {
				oathTknSrc := single.FromSupplierAsync(u.systemAccessTokenClient.GetGRPCAccessToken)
				return single.Map(oathTknSrc, gtools.OathAccessOption)
			},
			func() single.Single[grpc.DialOption] {
				return single.FromSupplierAsync(func() (result grpc.DialOption, err error) {
					return gtools.LoadTLSCredentialsOption(u.tlsConf.CACertPath())
				})
			},
		},
	)
	connectionSrc := single.MapWithError(optsSrc, func(opts []grpc.DialOption) (*grpc.ClientConn, error) {
		// Todo: load from config
		return gtools.CreateConnectionWithOptions("localhost:50051", opts...)
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

func NewUserService(systemAccessTokenClient SysAccessTokenClient, tlsConf conf.TLSConf) UserService {
	return &UserServiceImpl{systemAccessTokenClient: systemAccessTokenClient, tlsConf: tlsConf}
}
