package main

import (
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/businessrules"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/controllers"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/services/go/pkg/apperrors/errorservices"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/services/go/pkg/database/dbservices"
	"github.com/obenkenobi/cypher-log/services/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/services/go/pkg/logging"
	"github.com/obenkenobi/cypher-log/services/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/services/go/pkg/taskrunner"
	"github.com/obenkenobi/cypher-log/services/go/pkg/web/webservices"
)

func main() {
	environment.ReadEnvFiles(".env", "userservice.env") // Load env files
	logging.ConfigureGlobalLogging()                    // Configure logging

	// Dependency graph
	serverConf := conf.NewServerConf()
	auth0Conf := authconf.NewAuth0RouteSecurityConf()
	auth0ClientCredentialsConf := authconf.NewAuth0ClientCredentialsConf()
	mongoCOnf := conf.NewMongoConf()
	mongoHandler := dbservices.NewMongoHandler(mongoCOnf)
	userRepository := repositories.NewUserMongoRepository(mongoHandler)
	errorService := errorservices.NewErrorService()
	ginCtxService := webservices.NewGinCtxService(errorService)
	userBr := businessrules.NewUserBrImpl(mongoHandler, userRepository, errorService)
	authServerMgmtService := services.NewAuthServerMgmtService(auth0ClientCredentialsConf)
	userService := services.NewUserService(mongoHandler, userRepository, userBr, errorService, authServerMgmtService)
	authMiddleware := middlewares.NewAuthMiddleware(auth0Conf)
	userController := controllers.NewUserController(authMiddleware, userService, ginCtxService)
	appServer := webservices.NewServer(serverConf, userController)

	// Run tasks
	taskrunner.RunAndWait(func() { appServer.Run() })
}
