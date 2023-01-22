package servers

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/controllers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/commonservers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/controller"
	"net/http"
)

type AppServer interface {
	commonservers.CoreAppServer
}

type AppServerImpl struct {
	commonservers.CoreAppServer
}

//Todo: create environment configs, controllers and middlewares to properly configure the application.
//Todo: add csrf protection
//Todo: add proxies and middlewares to parse the auth token

func NewAppServerImpl(
	serverConf conf.ServerConf,
	tlsConf conf.TLSConf,
	authController controllers.AuthController,
	sessionMiddleware middlewares.SessionMiddleware,
	bearerAuthMiddleware middlewares.BearerAuthMiddleware,
) *AppServerImpl {
	beforeControllers := func(r *gin.Engine) {
		// Add gin engine configuration

		r.Use(sessionMiddleware.SessionHandler())
		r.Use(bearerAuthMiddleware.PassBearerTokenFromSession())

		r.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusPermanentRedirect, "/ui")
		})
		if environment.IsDevelopment() {
			r.LoadHTMLGlob("cmd/uiservice/resources/web/template/*")

		} else {
			r.Static("ui/", "cmd/uiservice/ClientApp/public")
		}

		if environment.IsDevelopment() {
			r.GET("/ui", func(ctx *gin.Context) {
				ctx.HTML(http.StatusOK, "home.html", nil)
			})
			r.GET("/ui/user", func(c *gin.Context) {
				c.HTML(http.StatusOK, "user.html", struct{}{})
			})
		}
	}

	controllersList := []controller.Controller{authController}

	afterControllers := func(r *gin.Engine) {
		/*Add gin engine configuration*/
	}

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
