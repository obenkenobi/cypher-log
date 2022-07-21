package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/controllers"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/services/go/pkg/middlewares"
	log "github.com/sirupsen/logrus"
)

type AppServer interface {
	Run()
}

type appserverImpl struct {
	serverConf conf.ServerConf
	router     *gin.Engine
}

func (s appserverImpl) Run() {
	err := s.router.Run(fmt.Sprintf(":%s", s.serverConf.GetPort()))
	if err != nil {
		log.Fatal("Server failed to run")
		return
	}
}

func BuildServer(userController controllers.UserController, serverConf conf.ServerConf,
	commonConf conf.CommonConf) AppServer {
	if commonConf.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
	server := appserverImpl{serverConf: serverConf, router: gin.New()}
	// Bind middlewares and routes
	middlewares.AddGlobalMiddleWares(server.router)
	userController.AddRoutes(server.router)
	return server
}
