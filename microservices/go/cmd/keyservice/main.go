package main

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/controllers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/listeners"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors/errorservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/database/dbservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/externalservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq/rmqservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security/securityservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/taskrunner"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/webservices"
)

func main() {
	environment.ReadEnvFiles(".env", "keyservice.env") // Load env files
	logger.ConfigureLoggerFromEnv()

	// Dependency graph
	serverConf := conf.NewServerConf()
	tlsConf := conf.NewTlsConf()
	httpclientConf := conf.NewHttpClientConf()
	grpcClientConf := conf.NewGrpcClientConf()
	auth0Conf := authconf.NewAuth0SecurityConf()
	rabbitMqConf := conf.NewRabbitMQConf()
	rabbitMqConnector := rmqservices.NewRabbitConnector(rabbitMqConf)
	defer rabbitMqConnector.Close()
	httpClientProvider := externalservices.NewHTTPClientProvider(httpclientConf)
	auth0SysAccessTokenClient := externalservices.NewAuth0SysAccessTokenClient(
		httpclientConf,
		auth0Conf,
		httpClientProvider,
	)
	coreGrpcConnProvider := externalservices.NewCoreGrpcConnProvider(auth0SysAccessTokenClient, tlsConf)
	extUserService := externalservices.NewExtUserService(coreGrpcConnProvider, grpcClientConf)
	mongoCOnf := conf.NewMongoConf()
	mongoHandler := dbservices.NewMongoHandler(mongoCOnf)
	userRepository := repositories.NewUserMongoRepository(mongoHandler)
	errorService := errorservices.NewErrorService()
	ginCtxService := webservices.NewGinCtxService(errorService)
	userService := services.NewUserService(userRepository, extUserService)

	// Add task dependencies
	var taskRunners []taskrunner.TaskRunner

	if environment.ActivateAppServer() { // Add app server
		apiAuth0JwtValidateService := securityservices.NewAPIAuth0JwtValidateService(auth0Conf)
		authMiddleware := middlewares.NewAuthMiddleware(apiAuth0JwtValidateService)
		userController := controllers.NewTestController(authMiddleware, userService, ginCtxService)
		appServer := webservices.NewAppServer(serverConf, tlsConf, userController)
		taskRunners = append(taskRunners, appServer)
	}
	if environment.ActivateRabbitMqListener() {
		rmqListener := listeners.NewRmqListener(rabbitMqConnector, userService)
		taskRunners = append(taskRunners, rmqListener)
	}

	// Run tasks
	taskrunner.RunAndWait(taskRunners...)

}
