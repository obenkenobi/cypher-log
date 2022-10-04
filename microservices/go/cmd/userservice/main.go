package main

import (
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
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/ginservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/grpcserveropts"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/rmqservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/securityservices"
)

func InitializeApp() *app.App {
	rabbitMQConfImpl := conf.NewRabbitMQConfImpl()
	rabbitMQPublisherImpl := rmqservices.NewRabbitPublisherImpl(rabbitMQConfImpl)
	serverConfImpl := conf.NewServerConfImpl()
	auth0RouteSecurityConfImpl := authconf.NewAuth0SecurityConfImpl()
	jwtValidateGrpcServiceImpl := securityservices.NewJwtValidateGrpcServiceImpl(auth0RouteSecurityConfImpl)
	authInterceptorCreatorImpl := grpcserveropts.NewAuthInterceptorCreatorImpl(jwtValidateGrpcServiceImpl)
	tlsConfImpl := conf.NewTlsConfImpl()
	credentialsOptionCreatorImpl := grpcserveropts.NewCredentialsOptionCreatorImpl(tlsConfImpl)
	userMessageServiceImpl := services.NewUserMessageServiceImpl(rabbitMQPublisherImpl)
	mongoConfImpl := conf.NewMongoConfImpl()
	mongoDBHandler := dshandlers.NewMongoDBHandler(mongoConfImpl)
	userRepositoryImpl := repositories.NewUserRepositoryImpl(mongoDBHandler)
	errorServiceImpl := sharedservices.NewErrorServiceImpl()
	userBrImpl := businessrules.NewUserBrImpl(mongoDBHandler, userRepositoryImpl, errorServiceImpl)
	authServerMgmtServiceImpl := services.NewAuthServerMgmtServiceImpl(auth0RouteSecurityConfImpl)
	userServiceImpl := services.NewUserServiceImpl(userMessageServiceImpl, mongoDBHandler, userRepositoryImpl, userBrImpl, errorServiceImpl, authServerMgmtServiceImpl)
	userServiceServerImpl := grpcservers.NewUserServiceServerImpl(userServiceImpl)
	grpcServerImpl := servers.NewGrpcServerImpl(serverConfImpl, authInterceptorCreatorImpl, credentialsOptionCreatorImpl, userServiceServerImpl)
	jwtValidateWebAppServiceImpl := securityservices.NewJwtValidateWebAppServiceImpl(auth0RouteSecurityConfImpl)
	authMiddlewareImpl := middlewares.NewAuthMiddlewareImpl(jwtValidateWebAppServiceImpl)
	ginCtxServiceImpl := ginservices.NewGinCtxServiceImpl(errorServiceImpl)
	userControllerImpl := controllers.NewUserControllerImpl(authMiddlewareImpl, userServiceImpl, ginCtxServiceImpl)
	appServerImpl := servers.NewAppServerImpl(serverConfImpl, tlsConfImpl, userControllerImpl)
	cronRunnerImpl := background.NewCronRunnerImpl(userServiceImpl)
	appApp := app.NewApp(rabbitMQPublisherImpl, grpcServerImpl, appServerImpl, cronRunnerImpl)
	return appApp
}

func main() {
	environment.ReadEnvFiles(".env", "userservice.env") // Load env files
	logger.ConfigureLoggerFromEnv()
	application := InitializeApp()
	application.Start()
}
