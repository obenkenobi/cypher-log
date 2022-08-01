package webservices

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf/environment"
	"github.com/obenkenobi/cypher-log/services/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/services/go/pkg/web"
	log "github.com/sirupsen/logrus"
)

// AppServer represents an interface that encapsulates the underlying web
// framework to act as a general application server that can be run.
type AppServer interface {
	// Run initiates the app server and listens now for underlying web requests.
	Run()
}

type appserverImpl struct {
	serverConf conf.ServerConf
	router     *gin.Engine
}

func (s appserverImpl) Run() {
	err := s.router.Run(fmt.Sprintf(":%s", s.serverConf.GetPort()))
	if err != nil {
		log.WithError(err).Fatal("Server failed to run")
		return
	}
}

// NewServer creates an app server that can be run by the server configuration and a list of controllers
func NewServer(serverConf conf.ServerConf, controllers ...web.Controller) AppServer {
	if environment.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
	server := appserverImpl{serverConf: serverConf, router: gin.New()}
	// Bind middlewares and routes
	middlewares.AddGlobalMiddleWares(server.router)
	for _, controller := range controllers {
		controller.AddRoutes(server.router)
	}
	return server
}
