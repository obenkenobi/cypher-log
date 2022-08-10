package main

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/controllers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors/errorservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/externalservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logging"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security/securityservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/taskrunner"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/webservices"
)

func main() {
	environment.ReadEnvFiles(".env", "keyservice.env") // Load env files
	logging.ConfigureGlobalLogging()                   // Configure logging

	// Dependency graph
	serverConf := conf.NewServerConf()
	tlsConf := conf.NewTlsConf()
	httpclientConf := conf.NewHttpClientConf()
	grpcClientConf := conf.NewGrpcClientConf()
	auth0Conf := authconf.NewAuth0SecurityConf()
	httpClientProvider := externalservices.NewHTTPClientProvider(httpclientConf)
	auth0SysAccessTokenClient := externalservices.NewAuth0SysAccessTokenClient(
		httpclientConf,
		auth0Conf,
		httpClientProvider,
	)
	coreGrpcConnProvider := externalservices.NewCoreGrpcConnProvider(auth0SysAccessTokenClient, tlsConf)
	userService := externalservices.NewExtUserService(coreGrpcConnProvider, grpcClientConf)
	//mongoCOnf := conf.NewMongoConf()
	//mongoHandler := dbservices.NewMongoHandler(mongoCOnf)
	errorService := errorservices.NewErrorService()
	ginCtxService := webservices.NewGinCtxService(errorService)

	// Add task dependencies
	var taskRunners []taskrunner.TaskRunner

	if environment.ActivateAppServer() { // Add app server
		apiAuth0JwtValidateService := securityservices.NewAPIAuth0JwtValidateService(auth0Conf)
		authMiddleware := middlewares.NewAuthMiddleware(apiAuth0JwtValidateService)
		userController := controllers.NewUserController(authMiddleware, userService, ginCtxService)
		appServer := webservices.NewAppServer(serverConf, tlsConf, userController)
		taskRunners = append(taskRunners, appServer)
	}

	// Run tasks
	taskrunner.RunAndWait(taskRunners...)

}
