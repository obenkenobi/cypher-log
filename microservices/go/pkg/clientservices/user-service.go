package clientservices

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
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
	grpcClientConf          conf.GrpcClientConf
	tlsConf                 conf.TLSConf
}

func (u UserServiceImpl) GetByAuthIdAsync(ctx context.Context, authId string) single.Single[userdtos.UserDto] {
	return single.MapToAsync[userdtos.UserDto](ctx, u.GetByAuthId(ctx, authId))
}
func (u UserServiceImpl) GetByAuthId(ctx context.Context, authId string) single.Single[userdtos.UserDto] {
	var dialOptSources []single.Single[grpc.DialOption]
	if environment.ActivateGRPCAuth() {
		oathTokenSrc := single.FromSupplierAsync(u.systemAccessTokenClient.GetGRPCAccessToken)
		oathOptSrc := single.Map(oathTokenSrc, gtools.OathAccessOption)
		tlsOptSrc := single.FromSupplierAsync(func() (grpc.DialOption, error) {
			return gtools.LoadTLSCredentialsOption(u.tlsConf.CACertPath(), environment.IsDevelopment())
		})
		dialOptSources = append(dialOptSources, oathOptSrc, tlsOptSrc)
	}
	optsSrc := gtools.CreateSingleWithDialOptions(dialOptSources)
	connectionSrc := single.MapWithError(optsSrc, func(opts []grpc.DialOption) (*grpc.ClientConn, error) {
		return gtools.CreateConnectionWithOptions(u.grpcClientConf.UserServiceAddress(), opts...)
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

func NewUserService(
	systemAccessTokenClient SysAccessTokenClient,
	tlsConf conf.TLSConf,
	grpcClientConf conf.GrpcClientConf,
) UserService {
	return &UserServiceImpl{
		systemAccessTokenClient: systemAccessTokenClient,
		tlsConf:                 tlsConf,
		grpcClientConf:          grpcClientConf,
	}
}
