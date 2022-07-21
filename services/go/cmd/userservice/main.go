package main

import (
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/controllers"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/services/go/pkg/middlewares"
)

// Environment variable names
const envVarKeyPort = "PORT"
const envVarKeyEnvironment = "ENVIRONMENT"
const envVarKeyAuth0IssuerUrl = "AUTH0_ISSUER_URL"
const envVarKeyAuth0Audience = "AUTH0_AUDIENCE"

func main() {
	// Dependency graph
	var envVarAccessor = conf.NewEnvVariableAccessor([]string{"userservice.env"})
	var serverConf = conf.NewServerConf(envVarAccessor, envVarKeyPort)
	var commonConf = conf.NewCommonConf(envVarAccessor, envVarKeyEnvironment)
	var auth0Conf = conf.NewAuth0Conf(envVarAccessor, envVarKeyAuth0IssuerUrl, envVarKeyAuth0Audience)
	var userService = services.NewUserService()
	var authMiddleware = middlewares.BuildAuthMiddleware(auth0Conf)
	var userController = controllers.NewUserController(authMiddleware, userService)
	var server = BuildServer(userController, serverConf, commonConf)

	// Run server
	server.Run()
}
