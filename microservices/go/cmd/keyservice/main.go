package main

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logging"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/taskrunner"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/webservices"
)

func main() {
	environment.ReadEnvFiles(".env", "keyservice.env") // Load env files
	logging.ConfigureGlobalLogging()                   // Configure logging

	// Dependency graph
	serverConf := conf.NewServerConf()
	//httpclientConf := conf.NewHttpClientConf()
	//auth0Conf := authconf.NewAuth0SecurityConf()
	//httpClientProvider := httpclient.NewClientProvider(httpclientConf)
	//auth0SysAccessTokenClient := httpclient.NewAuth0SysAccessTokenClient(
	//	httpclientConf,
	//	auth0Conf,
	//	httpClientProvider,
	//)
	//userService := clientservices.NewUserService(auth0SysAccessTokenClient)
	//mongoCOnf := conf.NewMongoConf()
	//mongoHandler := dbservices.NewMongoHandler(mongoCOnf)
	//errorService := errorservices.NewErrorService()
	//ginCtxService := webservices.NewGinCtxService(errorService)
	//authMiddleware := middlewares.NewAuthMiddleware(auth0Conf)
	appServer := webservices.NewAppServer(serverConf)

	// Run tasks
	taskrunner.RunAndWait(appServer)

}
