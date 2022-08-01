package main

import (
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf/environment"
	"github.com/obenkenobi/cypher-log/services/go/pkg/logging"
	"github.com/obenkenobi/cypher-log/services/go/pkg/web/webservices"
)

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
	environment.ReadEnvFiles(".env", "keyservice.env")                 // Load env files
	environment.SetEnvVarKeyForAppEnvironment(envVarKeyAppEnvironment) // Set app environment
	logging.ConfigureGlobalLogging()                                   // Configure logging

	// Dependency graph
	serverConf := conf.NewServerConf(envVarKeyPort)
	//auth0Conf := authconf.NewAuth0RouteSecurityConf(envVarKeyAuth0IssuerUrl, envVarKeyAuth0Audience)
	//auth0ClientCredentialsConf := authconf.NewAuth0ClientCredentialsConf(
	//	envVarKeyAuth0Domain,
	//	envVarKeyAuth0ClientId,
	//	envVarKeyAuth0ClientSecret,
	//	envVarKeyAuth0Audience,
	//)
	//mongoCOnf := conf.NewMongoConf(envVarKeyMongoUri, envVarMongoDBName, envVarMongoConnTimeoutMS)
	//mongoHandler := dbservices.NewMongoHandler(mongoCOnf)
	//errorService := errorservices.NewErrorService()
	//ginCtxService := webservices.NewGinCtxService(errorService)
	//authMiddleware := middlewares.NewAuthMiddleware(auth0Conf)
	appServer := webservices.NewServer(serverConf)

	// Run webservices
	appServer.Run()
}
