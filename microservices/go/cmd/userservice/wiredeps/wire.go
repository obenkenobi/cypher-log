//go:build wireinject
// +build wireinject

package wiredeps

import (
	"github.com/google/wire"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/app"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/background"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/businessrules"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/controllers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/grpcservers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/servers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userpb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/ginservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/grpcserveropts"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/rmqservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/securityservices"
)

func InitializeApp() *app.App {
	// wire_gen will not be imported but instead have InitializeApp be copied over to main.go
	wire.Build(
		conf.NewServerConfImpl,
		wire.Bind(new(conf.ServerConf), new(*conf.ServerConfImpl)),
		authconf.NewAuth0SecurityConfImpl,
		wire.Bind(new(authconf.Auth0SecurityConf), new(*authconf.Auth0RouteSecurityConfImpl)),
		conf.NewMongoConfImpl,
		wire.Bind(new(conf.MongoConf), new(*conf.MongoConfImpl)),
		conf.NewTlsConfImpl,
		wire.Bind(new(conf.TLSConf), new(*conf.TlsConfImpl)),
		conf.NewRabbitMQConfImpl,
		wire.Bind(new(conf.RabbitMQConf), new(*conf.RabbitMQConfImpl)),
		dshandlers.NewMongoDBHandler,
		wire.Bind(new(dshandlers.CrudDSHandler), new(*dshandlers.MongoDBHandler)),
		repositories.NewUserRepositoryImpl,
		wire.Bind(new(repositories.UserRepository), new(*repositories.UserRepositoryImpl)),
		rmqservices.NewRabbitPublisherImpl,
		wire.Bind(new(rmqservices.RabbitMQPublisher), new(*rmqservices.RabbitMQPublisherImpl)),
		sharedservices.NewErrorServiceImpl,
		wire.Bind(new(sharedservices.ErrorService), new(*sharedservices.ErrorServiceImpl)),
		businessrules.NewUserBrImpl,
		wire.Bind(new(businessrules.UserBr), new(*businessrules.UserBrImpl)),
		services.NewUserMessageServiceImpl,
		wire.Bind(new(services.UserMsgSendService), new(*services.UserMessageServiceImpl)),
		services.NewAuthServerMgmtServiceImpl,
		wire.Bind(new(services.AuthServerMgmtService), new(*services.AuthServerMgmtServiceImpl)),
		services.NewUserServiceImpl,
		wire.Bind(new(services.UserService), new(*services.UserServiceImpl)),
		ginservices.NewGinCtxServiceImpl,
		wire.Bind(new(ginservices.GinCtxService), new(*ginservices.GinCtxServiceImpl)),
		securityservices.NewJwtValidateWebAppServiceImpl,
		wire.Bind(new(securityservices.JwtValidateWebAppService), new(*securityservices.JwtValidateWebAppServiceImpl)),
		middlewares.NewAuthMiddlewareImpl,
		wire.Bind(new(middlewares.AuthMiddleware), new(*middlewares.AuthMiddlewareImpl)),
		controllers.NewUserControllerImpl,
		wire.Bind(new(controllers.UserController), new(*controllers.UserControllerImpl)),
		servers.NewAppServerImpl,
		wire.Bind(new(servers.AppServer), new(*servers.AppServerImpl)),
		securityservices.NewJwtValidateGrpcServiceImpl,
		wire.Bind(new(securityservices.JwtValidateGrpcService), new(*securityservices.JwtValidateGrpcServiceImpl)),
		grpcserveropts.NewAuthInterceptorCreatorImpl,
		wire.Bind(new(grpcserveropts.AuthInterceptorCreator), new(*grpcserveropts.AuthInterceptorCreatorImpl)),
		grpcserveropts.NewCredentialsOptionCreatorImpl,
		wire.Bind(new(grpcserveropts.CredentialsOptionCreator), new(*grpcserveropts.CredentialsOptionCreatorImpl)),
		grpcservers.NewUserServiceServerImpl,
		wire.Bind(new(userpb.UserServiceServer), new(*grpcservers.UserServiceServerImpl)),
		servers.NewGrpcServerImpl,
		wire.Bind(new(servers.GrpcServer), new(*servers.GrpcServerImpl)),
		background.NewCronRunnerImpl,
		wire.Bind(new(background.CronRunner), new(*background.CronRunnerImpl)),
		app.NewApp,
	)
	return &app.App{}
}
