package webservices

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/middlewares"
	log "github.com/sirupsen/logrus"
)

// AppServer represents an interface acts as a general application server that can be run.
type AppServer interface{ Server }

type appServerGinImpl struct {
	serverConf conf.ServerConf
	router     *gin.Engine
}

func (s appServerGinImpl) Run() {
	err := s.router.Run(fmt.Sprintf(":%s", s.serverConf.GetAppServerPort()))
	if err != nil {
		log.WithError(err).Fatal("AppServer failed to run")
		return
	}
}

// NewAppServer creates an app server that can be run by the server configuration and a list of controllers
func NewAppServer(serverConf conf.ServerConf, controllers ...Controller) AppServer {
	if environment.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
	server := appServerGinImpl{serverConf: serverConf, router: gin.New()}
	// Bind middlewares and routes
	middlewares.AddGlobalMiddleWares(server.router)
	for _, controller := range controllers {
		controller.AddRoutes(server.router)
	}
	return server
}
