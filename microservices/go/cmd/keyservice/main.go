package main

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/app"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/controllers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/listeners"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/commonservers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedrepos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/externalservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/ginservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/rmqservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/securityservices"
)

func main() {
	environment.ReadEnvFiles(".env", "keyservice.env") // Load env files
	logger.ConfigureLoggerFromEnv()

	// Dependency graph
	serverConf := conf.NewServerConfImpl()
	tlsConf := conf.NewTlsConfImpl()
	httpclientConf := conf.NewHttpClientConfImpl()
	grpcClientConf := conf.NewGrpcClientConfImpl()
	auth0Conf := authconf.NewAuth0SecurityConfImpl()
	rabbitMqConf := conf.NewRabbitMQConfImpl()
	rabbitConsumer := rmqservices.NewRabbitMQConsumerImpl(rabbitMqConf)
	httpClientProvider := externalservices.NewHTTPClientProviderImpl(httpclientConf)
	auth0SysAccessTokenClient := externalservices.NewSysAccessTokenClientAuth0Impl(
		httpclientConf,
		auth0Conf,
		httpClientProvider,
	)
	coreGrpcConnProvider := externalservices.NewCoreGrpcConnProviderImpl(auth0SysAccessTokenClient, tlsConf)
	extUserService := externalservices.NewExtUserServiceImpl(coreGrpcConnProvider, grpcClientConf)
	mongoCOnf := conf.NewMongoConfImpl()
	redisConf := conf.NewRedisConfImpl()
	mongoHandler := dshandlers.NewMongoDBHandler(mongoCOnf)
	redisHandler := dshandlers.NewRedisKeyValueTimedDBHandler(redisConf)
	userRepository := sharedrepos.NewUserRepositoryImpl(mongoHandler)
	userKeyRepository := repositories.NewUserKeyRepositoryImpl(mongoHandler)
	_ = repositories.NewAppSecretRepositoryImpl(redisHandler)
	_ = repositories.NewPrimaryAppSecretRefRepositoryImpl(redisHandler)
	errorService := sharedservices.NewErrorServiceImpl()
	ginCtxService := ginservices.NewGinCtxServiceImpl(errorService)
	userService := sharedservices.NewUserServiceImpl(userRepository, extUserService, errorService)
	userChangeEventService := services.NewUserChangeEventServiceImpl(userService, userKeyRepository, mongoHandler)
	ginRouterProvider := ginservices.NewGinEngineServiceImpl()
	apiAuth0JwtValidateService := securityservices.NewExternalOath2ValidateServiceAuth0Impl(auth0Conf)
	authMiddleware := middlewares.NewAuthMiddlewareImpl(apiAuth0JwtValidateService)
	_ = controllers.NewTestControllerImpl(ginRouterProvider, authMiddleware, userService, ginCtxService)
	appServer := commonservers.NewAppServerImpl(ginRouterProvider, serverConf, tlsConf)
	rmqListener := listeners.NewRmqListenerImpl(rabbitConsumer, userChangeEventService)

	application := app.NewApp(rabbitConsumer, appServer, rmqListener)
	application.Start()

}
