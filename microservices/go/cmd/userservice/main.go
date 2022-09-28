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
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/gservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/gservices/gserveroptions"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userpb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq/rmqservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security/securityservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/taskrunner"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/webservices"
	"google.golang.org/grpc"
)

func main() {
	environment.ReadEnvFiles(".env", "userservice.env") // Load env files

	// Main dependency graph
	serverConf := conf.NewServerConf()
	auth0Conf := authconf.NewAuth0SecurityConf()
	mongoCOnf := conf.NewMongoConf()
	mongoHandler := dbservices.NewMongoHandler(mongoCOnf)
	tlsConf := conf.NewTlsConf()
	rabbitMqConf := conf.NewRabbitMQConf()
	rabbitMqConnector := rmqservices.NewRabbitConnector(rabbitMqConf)
	defer rabbitMqConnector.Close()
	userMsgSendService := services.NewUserMessageServiceImpl(rabbitMqConnector)
	userRepository := repositories.NewUserMongoRepository(mongoHandler)
	errorService := errorservices.NewErrorService()
	userBr := businessrules.NewUserBrImpl(mongoHandler, userRepository, errorService)
	authServerMgmtService := services.NewAuthServerMgmtService(auth0Conf)
	userService := services.NewUserService(
		userMsgSendService,
		mongoHandler,
		userRepository,
		userBr,
		errorService,
		authServerMgmtService,
	)

	// Add task dependencies
	var taskRunners []taskrunner.TaskRunner

	if environment.ActivateAppServer() { // Add app server
		ginCtxService := webservices.NewGinCtxService(errorService)
		apiAuth0JwtValidateService := securityservices.NewAPIAuth0JwtValidateService(auth0Conf)
		authMiddleware := middlewares.NewAuthMiddleware(apiAuth0JwtValidateService)
		userController := controllers.NewUserController(authMiddleware, userService, ginCtxService)
		appServer := webservices.NewAppServer(serverConf, tlsConf, userController)
		taskRunners = append(taskRunners, appServer)
	}

	if environment.ActivateGrpcServer() { // Add GRPC server
		var grpcOpts []grpc.ServerOption
		if environment.ActivateGRPCAuth() {
			grpcAuth0JwtValidateService := securityservices.NewGrpcAuth0JwtValidateService(auth0Conf)
			authInterceptorCreator := gserveroptions.NewAuthInterceptorCreator(grpcAuth0JwtValidateService)
			credentialsOptionCreator := gserveroptions.NewCredentialsOptionCreator(tlsConf)
			grpcOpts = append(
				grpcOpts,
				authInterceptorCreator.CreateUnaryInterceptor(),
				credentialsOptionCreator.CreateCredentialsOption(),
			)
		}
		grpcServer := gservices.NewGrpcServer(
			serverConf,
			func(s *grpc.Server) {
				userpb.RegisterUserServiceServer(s, grpcservers.NewUserServiceServer(userService))
			},
			grpcOpts...,
		)
		taskRunners = append(taskRunners, grpcServer)
	}

	// Run taskRunners
	taskrunner.RunAndWait(taskRunners...)
}
