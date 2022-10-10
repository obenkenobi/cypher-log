//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/controllers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/listeners"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/servers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/services"
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
		conf.NewHttpClientConfImpl,
		wire.Bind(new(conf.HttpClientConf), new(*conf.HttpClientConfImpl)),
		conf.NewGrpcClientConfImpl,
		wire.Bind(new(conf.GrpcClientConf), new(*conf.GrpcClientConfImpl)),
		authconf.NewAuth0SecurityConfImpl,
		wire.Bind(new(authconf.Auth0SecurityConf), new(*authconf.Auth0RouteSecurityConfImpl)),
		conf.NewMongoConfImpl,
		wire.Bind(new(conf.MongoConf), new(*conf.MongoConfImpl)),
		conf.NewRabbitMQConfImpl,
		wire.Bind(new(conf.RabbitMQConf), new(*conf.RabbitMQConfImpl)),
		rmqservices.NewRabbitMQConsumerImpl,
		wire.Bind(new(rmqservices.RabbitMQConsumer), new(*rmqservices.RabbitMQConsumerImpl)),
		//appConf.NewKeyConfImpl,
		//wire.Bind(new(appConf.KeyConf), new(appConf.KeyConfImpl)),
		externalservices.NewHTTPClientProviderImpl,
		wire.Bind(new(externalservices.HttpClientProvider), new(*externalservices.HttpClientProviderImpl)),
		externalservices.NewSysAccessTokenClientAuth0Impl,
		wire.Bind(new(externalservices.SysAccessTokenClient), new(*externalservices.Auth0SysAccessTokenClient)),
		externalservices.NewCoreGrpcConnProviderImpl,
		wire.Bind(new(externalservices.CoreGrpcConnProvider), new(*externalservices.CoreGrpcConnProviderImpl)),
		externalservices.NewExtUserServiceImpl,
		wire.Bind(new(externalservices.ExtUserService), new(*externalservices.ExtUserServiceImpl)),
		//conf.NewRedisConfImpl,
		//wire.Bind(new(conf.RedisConf), new(*conf.RedisConfImpl)),
		dshandlers.NewMongoDBHandler,
		wire.Bind(new(dshandlers.CrudDSHandler), new(*dshandlers.MongoDBHandler)),
		//dshandlers.NewRedisDBHandler,
		//wire.Bind(new(dshandlers.KeyValueTimedDSHandler), new(*dshandlers.RedisDBHandler)),
		sharedrepos.NewUserRepositoryImpl,
		wire.Bind(new(sharedrepos.UserRepository), new(*sharedrepos.UserRepositoryImpl)),
		repositories.NewUserKeyRepositoryImpl,
		wire.Bind(new(repositories.UserKeyGeneratorRepository), new(*repositories.UserKeyGeneratorRepositoryImpl)),
		//repositories.NewUserKeySessionRepositoryImpl,
		//wire.Bind(new(repositories.UserKeySessionRepository), new(repositories.UserKeySessionRepositoryImpl)),
		//repositories.NewAppSecretRepositoryImpl,
		//wire.Bind(new(repositories.AppSecretRepository), new(*repositories.AppSecretRepositoryImpl)),
		//repositories.NewPrimaryAppSecretRefRepositoryImpl,
		//wire.Bind(new(repositories.PrimaryAppSecretRefRepository), new(*repositories.PrimaryAppSecretRefRepositoryImpl)),
		//rmqservices.NewRabbitPublisherImpl,
		//wire.Bind(new(rmqservices.RabbitMQPublisher), new(*rmqservices.RabbitMQPublisherImpl)),
		sharedservices.NewErrorServiceImpl,
		wire.Bind(new(sharedservices.ErrorService), new(*sharedservices.ErrorServiceImpl)),
		ginservices.NewGinCtxServiceImpl,
		wire.Bind(new(ginservices.GinCtxService), new(*ginservices.GinCtxServiceImpl)),
		sharedservices.NewUserServiceImpl,
		wire.Bind(new(sharedservices.UserService), new(*sharedservices.UserServiceImpl)),
		//businessrules.NewUserKeyBrImpl,
		//wire.Bind(new(businessrules.UserKeyBr), new(businessrules.UserKeyBrImpl)),
		services.NewUserChangeEventServiceImpl,
		wire.Bind(new(services.UserChangeEventService), new(*services.UserChangeEventServiceImpl)),
		//services.NewAppSecretServiceImpl,
		//wire.Bind(new(services.AppSecretService), new(services.AppSecretServiceImpl)),
		//NewUserKeyServiceImpl,
		//wire.Bind(new(services.UserKeyService), new(services.UserKeyServiceImpl)),
		securityservices.NewJwtValidateWebAppServiceImpl,
		wire.Bind(new(securityservices.JwtValidateWebAppService), new(*securityservices.JwtValidateWebAppServiceImpl)),
		middlewares.NewAuthMiddlewareImpl,
		wire.Bind(new(middlewares.AuthMiddleware), new(*middlewares.AuthMiddlewareImpl)),
		controllers.NewTestControllerImpl,
		wire.Bind(new(controllers.TestController), new(*controllers.TestControllerImpl)),
		servers.NewAppServerImpl,
		wire.Bind(new(servers.AppServer), new(*servers.AppServerImpl)),
		listeners.NewRmqListenerImpl,
		wire.Bind(new(listeners.RmqListener), new(*listeners.RmqListenerImpl)),
		NewApp,
	)
	return &App{}
}
