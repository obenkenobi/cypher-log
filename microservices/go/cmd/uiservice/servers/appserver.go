package servers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/controllers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/commonservers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
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

		r.Static("/public", "cmd/uiservice/resources/web/static")
		r.LoadHTMLGlob("cmd/uiservice/resources/web/template/*")

		r.Use(sessionMiddleware.SessionHandler())
		r.Use(bearerAuthMiddleware.PassBearerTokenFromSession())

		r.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "home.html", nil)
		})
		r.GET("/user", func(c *gin.Context) {
			session := sessions.Default(c)
			profile := session.Get(security.ProfileSessionKey)

			c.HTML(http.StatusOK, "user.html", profile)
		})
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
	if sessions.Default(ctx).Get("profile") == nil {
		ctx.Redirect(http.StatusSeeOther, "/")
	} else {
		ctx.Next()
	}
}
