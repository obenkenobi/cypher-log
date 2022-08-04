package main

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/businessrules"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/controllers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/grpcservers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors/errorservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/database/dbservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/grpcserveroptions"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userpb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logging"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security/securityservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/taskrunner"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/webservices"
	"google.golang.org/grpc"
)

func main() {
	environment.ReadEnvFiles(".env", "userservice.env") // Load env files
	logging.ConfigureGlobalLogging()                    // Configure logging

	// Main dependency graph
	serverConf := conf.NewServerConf()
	auth0Conf := authconf.NewAuth0RouteSecurityConf()
	mongoCOnf := conf.NewMongoConf()
	mongoHandler := dbservices.NewMongoHandler(mongoCOnf)
	userRepository := repositories.NewUserMongoRepository(mongoHandler)
	errorService := errorservices.NewErrorService()
	userBr := businessrules.NewUserBrImpl(mongoHandler, userRepository, errorService)
	authServerMgmtService := services.NewAuthServerMgmtService(auth0Conf)
	userService := services.NewUserService(mongoHandler, userRepository, userBr, errorService, authServerMgmtService)

	// Add task dependencies
	var tasks []func()

	if environment.ActivateAppServer() { // Add app server
		ginCtxService := webservices.NewGinCtxService(errorService)
		apiAuth0JwtValidateService := securityservices.NewAPIAuth0JwtValidateService(auth0Conf)
		authMiddleware := middlewares.NewAuthMiddleware(apiAuth0JwtValidateService)
		userController := controllers.NewUserController(authMiddleware, userService, ginCtxService)
		appServer := webservices.NewAppServer(serverConf, userController)
		tasks = append(tasks, appServer.Run)
	}

	if environment.ActivateGrpcServer() { // Add GRPC server
		grpcAuth0JwtValidateService := securityservices.NewGrpcAuth0JwtValidateService(auth0Conf)
		userServiceServer := grpcservers.NewUserServiceServer(userService)
		authInterceptorCreator := grpcserveroptions.NewAuthInterceptorCreator(grpcAuth0JwtValidateService)
		grpcServer := webservices.NewGrpcServer(
			serverConf,
			func(s *grpc.Server) {
				userpb.RegisterUserServiceServer(s, userServiceServer)
			},
			authInterceptorCreator.CreateUnaryInterceptor(),
		)
		tasks = append(tasks, grpcServer.Run)
	}

	// Run tasks
	taskrunner.RunAndWait(tasks...)
}
