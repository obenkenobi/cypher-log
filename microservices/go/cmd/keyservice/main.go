package main

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/controllers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/listeners"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/servers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedrepos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/externalservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/ginservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/rmqservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/securityservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/taskrunner"
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
	mongoHandler := dshandlers.NewMongoHandler(mongoCOnf)
	userRepository := sharedrepos.NewUserMongoRepository(mongoHandler)
	userKeyRepository := repositories.NewUserKeyRepository(mongoHandler)
	errorService := sharedservices.NewErrorService()
	ginCtxService := ginservices.NewGinCtxService(errorService)
	userService := sharedservices.NewUserService(userRepository, extUserService, errorService)
	userChangeEventService := services.NewUserChangeEventService(userService, userKeyRepository, mongoHandler)

	// Add task dependencies
	var taskRunners []taskrunner.TaskRunner

	if environment.ActivateAppServer() { // Add app server
		apiAuth0JwtValidateService := securityservices.NewAPIAuth0JwtValidateService(auth0Conf)
		authMiddleware := middlewares.NewAuthMiddleware(apiAuth0JwtValidateService)
		userController := controllers.NewTestController(authMiddleware, userService, ginCtxService)
		appServer := servers.NewAppServer(serverConf, tlsConf, userController)
		taskRunners = append(taskRunners, appServer)
	}
	if environment.ActivateRabbitMqListener() {
		rmqListener := listeners.NewRmqListener(rabbitMqConnector, userChangeEventService)
		taskRunners = append(taskRunners, rmqListener)
	}

	// Run tasks
	taskrunner.RunAndWait(taskRunners...)

}
