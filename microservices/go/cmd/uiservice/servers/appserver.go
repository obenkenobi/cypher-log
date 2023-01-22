package servers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/controllers"
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
) *AppServerImpl {

	gob.Register(map[string]interface{}{})
	store := cookie.NewStore([]byte("secret"))

	beforeControllers := func(r *gin.Engine) {
		// Add gin engine configuration

		r.Use(sessions.Sessions("auth-session", store))

		r.Static("/public", "cmd/uiservice/resources/web/static")

		r.LoadHTMLGlob("cmd/uiservice/resources/web/template/*")

		r.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "home.html", nil)
		})
		r.GET("/user", IsAuthenticated, func(c *gin.Context) {
			session := sessions.Default(c)
			profile := session.Get(security.ProfileSessionKey)

			c.HTML(http.StatusOK, "user.html", profile)
		})
	}
	controllersList := []controller.Controller{authController}
	afterControllers := func(r *gin.Engine) {
		// Add gin engine configuration
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
	if sessions.Default(ctx).Get("profile") == nil {
		ctx.Redirect(http.StatusSeeOther, "/")
	} else {
		ctx.Next()
	}
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
