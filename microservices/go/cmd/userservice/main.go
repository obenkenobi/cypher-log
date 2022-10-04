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
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/commonservers"
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

func main() {
	environment.ReadEnvFiles(".env", "userservice.env") // Load env files
	logger.ConfigureLoggerFromEnv()

	// Main dependency graph
	serverConf := conf.NewServerConfImpl()
	auth0Conf := authconf.NewAuth0SecurityConfImpl()
	mongoCOnf := conf.NewMongoConfImpl()
	tlsConf := conf.NewTlsConfImpl()
	rabbitMqConf := conf.NewRabbitMQConfImpl()
	mongoHandler := dshandlers.NewMongoDBHandler(mongoCOnf)
	rabbitPublisher := rmqservices.NewRabbitPublisherImpl(rabbitMqConf)
	userMsgSendService := services.NewUserMessageServiceImpl(rabbitPublisher)
	userRepository := repositories.NewUserRepositoryImpl(mongoHandler)
	errorService := sharedservices.NewErrorServiceImpl()
	userBr := businessrules.NewUserBrImpl(mongoHandler, userRepository, errorService)
	authServerMgmtService := services.NewAuthServerMgmtServiceImpl(auth0Conf)
	userService := services.NewUserServiceImpl(
		userMsgSendService,
		mongoHandler,
		userRepository,
		userBr,
		errorService,
		authServerMgmtService,
	)
	ginRouterProvider := ginservices.NewGinEngineServiceImpl()
	ginCtxService := ginservices.NewGinCtxServiceImpl(errorService)
	apiAuth0JwtValidateService := securityservices.NewExternalOath2ValidateServiceAuth0Impl(auth0Conf)
	authMiddleware := middlewares.NewAuthMiddlewareImpl(apiAuth0JwtValidateService)
	_ = controllers.NewUserControllerImpl(ginRouterProvider, authMiddleware, userService, ginCtxService)
	appServer := commonservers.NewAppServerImpl(ginRouterProvider, serverConf, tlsConf)

	grpcAuth0JwtValidateService := securityservices.NewExternalOath2ValidateServiceImpl(auth0Conf)
	authInterceptorCreator := grpcserveropts.NewAuthInterceptorCreatorImpl(grpcAuth0JwtValidateService)
	credentialsOptionCreator := grpcserveropts.NewCredentialsOptionCreatorImpl(tlsConf)
	userServiceServer := grpcservers.NewUserServiceServerImpl(userService)
	grpcServer := servers.NewGrpcServerImpl(
		serverConf,
		authInterceptorCreator,
		credentialsOptionCreator,
		userServiceServer,
	)
	cronRunner := background.NewCronRunnerImpl(userService)

	application := app.NewApp(rabbitPublisher, grpcServer, appServer, cronRunner)
	application.Start()
}
