package main

import (
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/businessrules"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/controllers"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dbaccess"
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
	// Dependency graph
	var envVarReader = conf.NewEnvVariableAccessor([]string{"userservice.env"})
	var serverConf = conf.NewServerConf(envVarReader, envVarKeyPort)
	var commonConf = conf.NewCommonConf(envVarReader, envVarKeyEnvironment)
	var auth0Conf = conf.NewAuth0Conf(envVarReader, envVarKeyAuth0IssuerUrl, envVarKeyAuth0Audience)
	var mongoCOnf = conf.NewMongoConf(envVarReader, envVarKeyMongoUri, envVarMongoDBName, envVarMongoConnTimeoutMS)
	var mongoClient = dbaccess.BuildMongoClient(mongoCOnf)
	var transactionRunner = dbaccess.NewTransactionRunnerMongo(mongoClient)
	var userRepository = repositories.NewUserMongoRepository(mongoClient)
	var userBr = businessrules.NewUserBrImpl(mongoClient, userRepository)
	var userService = services.NewUserService(mongoClient, transactionRunner, userRepository, userBr)
	var authMiddleware = middlewares.BuildAuthMiddleware(auth0Conf)
	var userController = controllers.NewUserController(authMiddleware, userService)
	var server = BuildServer(userController, serverConf, commonConf)

	// Run server
	server.Run()
}
