package main

import (
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/businessrules"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/controllers"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/services/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf/environment"
	"github.com/obenkenobi/cypher-log/services/go/pkg/database"
	"github.com/obenkenobi/cypher-log/services/go/pkg/framework/ginextensions"
	"github.com/obenkenobi/cypher-log/services/go/pkg/logging"
	"github.com/obenkenobi/cypher-log/services/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/services/go/pkg/server"
)

// Environment variable names
const envVarKeyPort = "PORT"
const envVarKeyAppEnvironment = "ENVIRONMENT"
const envVarKeyAuth0IssuerUrl = "AUTH0_ISSUER_URL"
const envVarKeyAuth0Audience = "AUTH0_AUDIENCE"
const envVarKeyMongoUri = "MONGO_URI"
const envVarMongoDBName = "MONGO_DB_NAME"
const envVarMongoConnTimeoutMS = "MONGO_CONNECTION_TIMEOUT_MS"

func main() {
	environment.ReadEnvFiles("userservice.env")                        // Load env files
	environment.SetEnvVarKeyForAppEnvironment(envVarKeyAppEnvironment) // Set app environment
	logging.ConfigureGlobalLogging()                                   // Configure logging
	// Dependency graph
	var serverConf = conf.NewServerConf(envVarKeyPort)
	var auth0Conf = authconf.NewAuth0Conf(envVarKeyAuth0IssuerUrl, envVarKeyAuth0Audience)
	var mongoCOnf = conf.NewMongoConf(envVarKeyMongoUri, envVarMongoDBName, envVarMongoConnTimeoutMS)
	var mongoHandler = database.BuildMongoHandler(mongoCOnf)
	var userRepository = repositories.NewUserMongoRepository(mongoHandler)
	var errorService = apperrors.NewErrorServiceImpl()
	var ginCtxService = ginextensions.NewGinWrapperService(errorService)
	var userBr = businessrules.NewUserBrImpl(mongoHandler, userRepository, errorService)
	var userService = services.NewUserService(mongoHandler, userRepository, userBr, errorService)
	var authMiddleware = middlewares.BuildAuthMiddleware(auth0Conf)
	var userController = controllers.NewUserController(authMiddleware, userService, ginCtxService)
	var appServer = server.BuildServer(serverConf, userController)

	// Run server
	appServer.Run()
}
