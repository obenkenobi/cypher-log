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
//Todo: add proxies to apis

func NewAppServerImpl(
	serverConf conf.ServerConf,
	tlsConf conf.TLSConf,
	authController controllers.AuthController,
	sessionMiddleware middlewares.SessionMiddleware,
	bearerAuthMiddleware middlewares.BearerAuthMiddleware,
	uiProviderMiddleware middlewares.UiProviderMiddleware,
) *AppServerImpl {
	beforeControllers := func(r *gin.Engine) {
		// Add gin engine configuration
		r.Use(sessionMiddleware.SessionHandler())
		r.Use(bearerAuthMiddleware.PassBearerTokenFromSession())
		uiProviderMiddleware.ProvideUI(r)
	}
	controllersList := []controller.Controller{authController}
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

func IsAuthenticated(ctx *gin.Context) {
}
