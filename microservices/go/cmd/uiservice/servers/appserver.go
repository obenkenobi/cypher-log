package servers

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/controllers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/commonservers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/lifecycle"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/controller"
)

type AppServer interface {
	commonservers.CoreAppServer
}

type AppServerImpl struct {
	commonservers.CoreAppServer
}

func NewAppServerImpl(
	serverConf conf.ServerConf,
	tlsConf conf.TLSConf,
	sessionMiddleware middlewares.SessionMiddleware,
	uiProviderMiddleware middlewares.UiProviderMiddleware,
	authController controllers.AuthController,
	csrfController controllers.CsrfController,
	gatewayController controllers.GatewayController,
) *AppServerImpl {
	if !environment.ActivateAppServer() {
		// App server is deactivated, ran via the lifecycle package,
		// and is a root-child dependency so a nil is returned
		return nil
	}
	beforeControllers := func(r *gin.Engine) {
		// Add gin engine configuration
		uiProviderMiddleware.ProvideUI(r)
		r.Use(sessionMiddleware.SessionHandler())
		r.Use(sessionMiddleware.CsrfHandler())
	}
	controllersList := []controller.Controller{authController, csrfController, gatewayController}
	afterControllers := func(r *gin.Engine) { /*Add gin engine configuration*/ }

	coreAppServer := commonservers.NewCoreAppServerWithHooksImpl(
		serverConf,
		tlsConf,
		beforeControllers,
		controllersList,
		afterControllers,
	)
	a := &AppServerImpl{CoreAppServer: coreAppServer}
	lifecycle.RegisterTaskRunner(a)
	return a
}
