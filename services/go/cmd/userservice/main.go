package main

import (
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/businessrules"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/controllers"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/services/go/pkg/apperrors/errorservices"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf/environment"
	"github.com/obenkenobi/cypher-log/services/go/pkg/database"
	"github.com/obenkenobi/cypher-log/services/go/pkg/logging"
	"github.com/obenkenobi/cypher-log/services/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/services/go/pkg/web/webservices"
)

// Environment variable names
const envVarKeyPort = "PORT"
const envVarKeyAppEnvironment = "ENVIRONMENT"
const envVarKeyAuth0IssuerUrl = "AUTH0_ISSUER_URL"
const envVarKeyAuth0Audience = "AUTH0_AUDIENCE"
const envVarKeyAuth0Domain = "AUTH0_DOMAIN"
const envVarKeyAuth0ClientId = "AUTH0_CLIENT_ID"
const envVarKeyAuth0ClientSecret = "AUTH0_CLIENT_SECRET"
const envVarKeyMongoUri = "MONGO_URI"
const envVarMongoDBName = "MONGO_DB_NAME"
const envVarMongoConnTimeoutMS = "MONGO_CONNECTION_TIMEOUT_MS"

func main() {
	environment.ReadEnvFiles("userservice.env")                        // Load env files
	environment.SetEnvVarKeyForAppEnvironment(envVarKeyAppEnvironment) // Set app environment
	logging.ConfigureGlobalLogging()                                   // Configure logging
	// Dependency graph
	serverConf := conf.NewServerConf(envVarKeyPort)
	auth0Conf := authconf.NewAuth0RouteSecurityConf(envVarKeyAuth0IssuerUrl, envVarKeyAuth0Audience)
	auth0ClientCredentialsConf := authconf.NewAuth0ClientCredentialsConf(
		envVarKeyAuth0Domain,
		envVarKeyAuth0ClientId,
		envVarKeyAuth0ClientSecret,
		envVarKeyAuth0Audience,
	)
	mongoCOnf := conf.NewMongoConf(envVarKeyMongoUri, envVarMongoDBName, envVarMongoConnTimeoutMS)
	mongoHandler := database.BuildMongoHandler(mongoCOnf)
	userRepository := repositories.NewUserMongoRepository(mongoHandler)
	errorService := errorservices.NewErrorService()
	ginCtxService := webservices.NewGinWrapperService(errorService)
	userBr := businessrules.NewUserBrImpl(mongoHandler, userRepository, errorService)
	authServerMgmtService := services.NewAuthServerMgmtService(auth0ClientCredentialsConf)
	userService := services.NewUserService(mongoHandler, userRepository, userBr, errorService, authServerMgmtService)
	authMiddleware := middlewares.BuildAuthMiddleware(auth0Conf)
	userController := controllers.NewUserController(authMiddleware, userService, ginCtxService)
	appServer := webservices.BuildServer(serverConf, userController)

	// Run webservices
	appServer.Run()
}
