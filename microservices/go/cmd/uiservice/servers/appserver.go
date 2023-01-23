package servers

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/controllers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/commonservers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/controller"
)

type AppServer interface {
	commonservers.CoreAppServer
}

type AppServerImpl struct {
	commonservers.CoreAppServer
}

//Todo: add csrf protection

func NewAppServerImpl(
	serverConf conf.ServerConf,
	tlsConf conf.TLSConf,
	sessionMiddleware middlewares.SessionMiddleware,
	uiProviderMiddleware middlewares.UiProviderMiddleware,
	userKeyMiddleware middlewares.UserKeyMiddleware,
	authController controllers.AuthController,
	gatewayController controllers.GatewayController,
) *AppServerImpl {
	beforeControllers := func(r *gin.Engine) {
		// Add gin engine configuration
		uiProviderMiddleware.ProvideUI(r)
		r.Use(sessionMiddleware.SessionHandler())
		r.Use(userKeyMiddleware.UserKeySession())
	}
	controllersList := []controller.Controller{authController, gatewayController}
	afterControllers := func(r *gin.Engine) { /*Add gin engine configuration*/ }

	coreAppServer := commonservers.NewCoreAppServerWithHooksImpl(
		serverConf,
		tlsConf,
		beforeControllers,
		controllersList,
		afterControllers,
	)
	return &AppServerImpl{CoreAppServer: coreAppServer}
}
