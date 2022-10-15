//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/controllers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/listeners"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/servers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedrepos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/externalservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/ginservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/rmqservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/securityservices"
)

func InitializeApp() *App {
	wire.Build(
		conf.NewServerConfImpl,
		wire.Bind(new(conf.ServerConf), new(*conf.ServerConfImpl)),
		conf.NewTlsConfImpl,
		wire.Bind(new(conf.TLSConf), new(*conf.TlsConfImpl)),
		conf.NewMongoConfImpl,
		wire.Bind(new(conf.MongoConf), new(*conf.MongoConfImpl)),
		authconf.NewAuth0SecurityConfImpl,
		wire.Bind(new(authconf.Auth0SecurityConf), new(*authconf.Auth0RouteSecurityConfImpl)),
		conf.NewGrpcClientConfImpl,
		wire.Bind(new(conf.GrpcClientConf), new(*conf.GrpcClientConfImpl)),
		conf.NewHttpClientConfImpl,
		wire.Bind(new(conf.HttpClientConf), new(*conf.HttpClientConfImpl)),
		conf.NewRabbitMQConfImpl,
		wire.Bind(new(conf.RabbitMQConf), new(*conf.RabbitMQConfImpl)),
		rmqservices.NewRabbitMQConsumerImpl,
		wire.Bind(new(rmqservices.RabbitMQConsumer), new(*rmqservices.RabbitMQConsumerImpl)),
		externalservices.NewSysAccessTokenClientAuth0Impl,
		wire.Bind(new(externalservices.SysAccessTokenClient), new(*externalservices.Auth0SysAccessTokenClient)),
		externalservices.NewHTTPClientProviderImpl,
		wire.Bind(new(externalservices.HttpClientProvider), new(*externalservices.HttpClientProviderImpl)),
		externalservices.NewExtUserServiceImpl,
		wire.Bind(new(externalservices.ExtUserService), new(*externalservices.ExtUserServiceImpl)),
		externalservices.NewCoreGrpcConnProviderImpl,
		wire.Bind(new(externalservices.CoreGrpcConnProvider), new(*externalservices.CoreGrpcConnProviderImpl)),
		dshandlers.NewMongoDBHandler,
		wire.Bind(new(dshandlers.CrudDSHandler), new(*dshandlers.MongoDBHandler)),
		sharedrepos.NewUserRepositoryImpl,
		wire.Bind(new(sharedrepos.UserRepository), new(*sharedrepos.UserRepositoryImpl)),
		ginservices.NewGinCtxServiceImpl,
		wire.Bind(new(ginservices.GinCtxService), new(*ginservices.GinCtxServiceImpl)),
		sharedservices.NewUserServiceImpl,
		wire.Bind(new(sharedservices.UserService), new(*sharedservices.UserServiceImpl)),
		sharedservices.NewErrorServiceImpl,
		wire.Bind(new(sharedservices.ErrorService), new(*sharedservices.ErrorServiceImpl)),
		repositories.NewNoteRepositoryImpl,
		wire.Bind(new(repositories.NoteRepository), new(*repositories.NoteRepositoryImpl)),
		services.NewUserChangeEventServiceImpl,
		wire.Bind(new(services.UserChangeEventService), new(*services.UserChangeEventServiceImpl)),
		services.NewNoteServiceImpl,
		wire.Bind(new(services.NoteService), new(*services.NoteServiceImpl)),
		securityservices.NewJwtValidateWebAppServiceImpl,
		wire.Bind(new(securityservices.JwtValidateWebAppService), new(*securityservices.JwtValidateWebAppServiceImpl)),
		middlewares.NewAuthMiddlewareImpl,
		wire.Bind(new(middlewares.AuthMiddleware), new(*middlewares.AuthMiddlewareImpl)),
		controllers.NewNoteControllerImpl,
		wire.Bind(new(controllers.NoteController), new(*controllers.NoteControllerImpl)),
		servers.NewAppServerImpl,
		wire.Bind(new(servers.AppServer), new(*servers.AppServerImpl)),
		listeners.NewRmqListenerImpl,
		wire.Bind(new(listeners.RmqListener), new(*listeners.RmqListenerImpl)),
		NewApp)
	return &App{}
}
