package main

import (
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/businessrules"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/controllers"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/services/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/services/go/pkg/database"
	"github.com/obenkenobi/cypher-log/services/go/pkg/framework/ginextensions"
	"github.com/obenkenobi/cypher-log/services/go/pkg/logging"
	"github.com/obenkenobi/cypher-log/services/go/pkg/middlewares"
)

// Environment variable names
const envVarKeyPort = "PORT"
const envVarKeyEnvironment = "ENVIRONMENT"
const envVarKeyAuth0IssuerUrl = "AUTH0_ISSUER_URL"
const envVarKeyAuth0Audience = "AUTH0_AUDIENCE"
const envVarKeyMongoUri = "MONGO_URI"
const envVarMongoDBName = "MONGO_DB_NAME"
const envVarMongoConnTimeoutMS = "MONGO_CONNECTION_TIMEOUT_MS"

func main() {
	logging.ConfigTextLogging()

	// Dependency graph
	var envVarReader = conf.NewEnvVariableReader([]string{"userservice.env"})
	var serverConf = conf.NewServerConf(envVarReader, envVarKeyPort)
	var commonConf = conf.NewCommonConf(envVarReader, envVarKeyEnvironment)
	var auth0Conf = conf.NewAuth0Conf(envVarReader, envVarKeyAuth0IssuerUrl, envVarKeyAuth0Audience)
	var mongoCOnf = conf.NewMongoConf(envVarReader, envVarKeyMongoUri, envVarMongoDBName, envVarMongoConnTimeoutMS)
	var mongoHandler = database.BuildMongoHandler(mongoCOnf)
	var userRepository = repositories.NewUserMongoRepository(mongoHandler)
	var errorService = apperrors.NewErrorServiceImpl()
	var ginWrapperService = ginextensions.NewGinWrapperService(errorService)
	var userBr = businessrules.NewUserBrImpl(mongoHandler, userRepository, errorService)
	var userService = services.NewUserService(mongoHandler, userRepository, userBr, errorService)
	var authMiddleware = middlewares.BuildAuthMiddleware(auth0Conf)
	var userController = controllers.NewUserController(authMiddleware, userService, ginWrapperService)
	var server = BuildServer(userController, serverConf, commonConf)

	// Run server
	server.Run()
}
