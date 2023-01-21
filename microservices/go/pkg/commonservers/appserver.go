package commonservers

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

// CoreAppServer represents an interface acts as a general application server that can be run as a task.
type CoreAppServer interface{ taskrunner.TaskRunner }

type CoreAppServerImpl struct {
	router     *gin.Engine
	serverConf conf.ServerConf
	tlsConf    conf.TLSConf
}

func (s CoreAppServerImpl) Run() {
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
		logger.Log.WithError(err).Fatal("CoreAppServer failed to run")
		return
	}
}

func (s CoreAppServerImpl) getFormattedPort() string {
	return fmt.Sprintf(":%s", s.serverConf.GetAppServerPort())
}

// NewCoreAppServerImpl creates an app server that can be run by the server configuration and a list of controllers
func NewCoreAppServerImpl(
	serverConf conf.ServerConf,
	tlsConf conf.TLSConf,
	controllers ...controller.Controller,
) *CoreAppServerImpl {
	return NewCoreAppServerWithHooksImpl(
		serverConf,
		tlsConf,
		func(r *gin.Engine) {},
		controllers,
		func(r *gin.Engine) {},
	)
}

// NewCoreAppServerWithHooksImpl creates an app server that can be run by the
// server configuration and a list of controllers. The beforeControllers and
// afterControllers hooks add additional gin engine configuration before and
// after controllers are added respectively.
func NewCoreAppServerWithHooksImpl(
	serverConf conf.ServerConf,
	tlsConf conf.TLSConf,
	beforeControllers func(r *gin.Engine),
	controllers []controller.Controller,
	afterControllers func(r *gin.Engine),
) *CoreAppServerImpl {
	if environment.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	middlewares.AddGlobalMiddleWares(r)
	beforeControllers(r)
	for _, c := range controllers {
		c.AddRoutes(r)
	}
	afterControllers(r)
	server := &CoreAppServerImpl{serverConf: serverConf, tlsConf: tlsConf, router: r}
	return server
}
