package servers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/taskrunner"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/controller"
)

// AppServer represents an interface acts as a general application server that can be run as a task.
type AppServer interface{ taskrunner.TaskRunner }

type appServerGinImpl struct {
	serverConf conf.ServerConf
	tlsConf    conf.TLSConf
	router     *gin.Engine
}

func (s appServerGinImpl) Run() {
	var err error
	if environment.ActivateAppServerTLS() {
		err = s.router.RunTLS(
			s.getFormattedPort(),
			s.tlsConf.ServerCertPath(),
			s.tlsConf.ServerKeyPath(),
		)
	} else {
		err = s.router.Run(s.getFormattedPort())
	}
	if err != nil {
		logger.Log.WithError(err).Fatal("AppServer failed to run")
		return
	}
}

func (s appServerGinImpl) getFormattedPort() string {
	return fmt.Sprintf(":%s", s.serverConf.GetAppServerPort())
}

// NewAppServer creates an app server that can be run by the server configuration and a list of controllers
func NewAppServer(serverConf conf.ServerConf, tlsConf conf.TLSConf, controllers ...controller.Controller) AppServer {
	if environment.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
	server := appServerGinImpl{serverConf: serverConf, tlsConf: tlsConf, router: gin.New()}
	// Bind middlewares and routes
	middlewares.AddGlobalMiddleWares(server.router)
	for _, controller := range controllers {
		controller.AddRoutes(server.router)
	}
	return server
}
