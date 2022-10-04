package commonservers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/ginservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/taskrunner"
)

// AppServer represents an interface acts as a general application server that can be run as a task.
type AppServer interface{ taskrunner.TaskRunner }

type AppServerImpl struct {
	ginRouterProvider ginservices.GinRouterProvider
	serverConf        conf.ServerConf
	tlsConf           conf.TLSConf
}

func (s AppServerImpl) Run() {
	s.ginRouterProvider.AccessRouter(func(r *gin.Engine) {
		var err error
		if environment.ActivateAppServerTLS() {
			err = r.RunTLS(
				s.getFormattedPort(),
				s.tlsConf.ServerCertPath(),
				s.tlsConf.ServerKeyPath(),
			)
		} else {
			err = r.Run(s.getFormattedPort())
		}
		if err != nil {
			logger.Log.WithError(err).Fatal("AppServer failed to run")
			return
		}
	})

}

func (s AppServerImpl) getFormattedPort() string {
	return fmt.Sprintf(":%s", s.serverConf.GetAppServerPort())
}

// NewAppServerImpl creates an app server that can be run by the server configuration and a list of controllers
func NewAppServerImpl(
	ginRouterProvider ginservices.GinRouterProvider,
	serverConf conf.ServerConf,
	tlsConf conf.TLSConf,
) *AppServerImpl {
	ginRouterProvider.AccessRouter(middlewares.AddGlobalMiddleWares)
	server := &AppServerImpl{serverConf: serverConf, tlsConf: tlsConf, ginRouterProvider: ginRouterProvider}
	return server
}
