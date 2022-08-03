package main

import (
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/services/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/services/go/pkg/logging"
	"github.com/obenkenobi/cypher-log/services/go/pkg/taskrunner"
	"github.com/obenkenobi/cypher-log/services/go/pkg/web/webservices"
)

func main() {
	environment.ReadEnvFiles(".env", "keyservice.env") // Load env files
	logging.ConfigureGlobalLogging()                   // Configure logging

	// Dependency graph
	serverConf := conf.NewServerConf()
	//auth0Conf := authconf.NewAuth0RouteSecurityConf()
	//auth0ClientCredentialsConf := authconf.NewAuth0ClientCredentialsConf()
	//mongoCOnf := conf.NewMongoConf()
	//mongoHandler := dbservices.NewMongoHandler(mongoCOnf)
	//errorService := errorservices.NewErrorService()
	//ginCtxService := webservices.NewGinCtxService(errorService)
	//authMiddleware := middlewares.NewAuthMiddleware(auth0Conf)
	appServer := webservices.NewServer(serverConf)

	// Run tasks
	taskrunner.RunAndWait(func() { appServer.Run() })

}
