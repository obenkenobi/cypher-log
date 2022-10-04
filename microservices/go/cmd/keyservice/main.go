package main

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/app"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/controllers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/listeners"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/servers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/services"
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

func InitializeApp() *app.App {
	rabbitMQConfImpl := conf.NewRabbitMQConfImpl()
	rabbitMQConsumerImpl := rmqservices.NewRabbitMQConsumerImpl(rabbitMQConfImpl)
	serverConfImpl := conf.NewServerConfImpl()
	tlsConfImpl := conf.NewTlsConfImpl()
	auth0RouteSecurityConfImpl := authconf.NewAuth0SecurityConfImpl()
	jwtValidateWebAppServiceImpl := securityservices.NewJwtValidateWebAppServiceImpl(auth0RouteSecurityConfImpl)
	authMiddlewareImpl := middlewares.NewAuthMiddlewareImpl(jwtValidateWebAppServiceImpl)
	mongoConfImpl := conf.NewMongoConfImpl()
	mongoDBHandler := dshandlers.NewMongoDBHandler(mongoConfImpl)
	userRepositoryImpl := sharedrepos.NewUserRepositoryImpl(mongoDBHandler)
	httpClientConfImpl := conf.NewHttpClientConfImpl()
	httpClientProviderImpl := externalservices.NewHTTPClientProviderImpl(httpClientConfImpl)
	auth0SysAccessTokenClient := externalservices.NewSysAccessTokenClientAuth0Impl(httpClientConfImpl, auth0RouteSecurityConfImpl, httpClientProviderImpl)
	coreGrpcConnProviderImpl := externalservices.NewCoreGrpcConnProviderImpl(auth0SysAccessTokenClient, tlsConfImpl)
	grpcClientConfImpl := conf.NewGrpcClientConfImpl()
	extUserServiceImpl := externalservices.NewExtUserServiceImpl(coreGrpcConnProviderImpl, grpcClientConfImpl)
	errorServiceImpl := sharedservices.NewErrorServiceImpl()
	userServiceImpl := sharedservices.NewUserServiceImpl(userRepositoryImpl, extUserServiceImpl, errorServiceImpl)
	ginCtxServiceImpl := ginservices.NewGinCtxServiceImpl(errorServiceImpl)
	testControllerImpl := controllers.NewTestControllerImpl(authMiddlewareImpl, userServiceImpl, ginCtxServiceImpl)
	appServerImpl := servers.NewAppServerImpl(serverConfImpl, tlsConfImpl, testControllerImpl)
	userKeyRepositoryImpl := repositories.NewUserKeyRepositoryImpl(mongoDBHandler)
	userChangeEventServiceImpl := services.NewUserChangeEventServiceImpl(userServiceImpl, userKeyRepositoryImpl, mongoDBHandler)
	rmqListenerImpl := listeners.NewRmqListenerImpl(rabbitMQConsumerImpl, userChangeEventServiceImpl)
	appApp := app.NewApp(rabbitMQConsumerImpl, appServerImpl, rmqListenerImpl)
	return appApp
}

func main() {
	environment.ReadEnvFiles(".env", "keyservice.env") // Load env files
	logger.ConfigureLoggerFromEnv()
	application := InitializeApp()
	application.Start()

}
