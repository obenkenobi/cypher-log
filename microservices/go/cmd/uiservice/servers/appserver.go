package servers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/commonservers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/controller"
	"net/http"
	"net/url"
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
	authenticatorService services.AuthenticatorService,
	authConf authconf.Auth0SecurityConf,
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

		// Login
		r.GET("/login", func(ctx *gin.Context) {
			state, err := generateRandomState()
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}

			// Save the state inside the session.
			session := sessions.Default(ctx)
			session.Set("state", state)
			if err := session.Save(); err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
			ctx.Redirect(http.StatusTemporaryRedirect, authenticatorService.GetOath2Config().AuthCodeURL(state))
		})

		// Callback
		r.GET("/callback", func(ctx *gin.Context) {
			session := sessions.Default(ctx)
			if ctx.Query("state") != session.Get("state") {
				ctx.String(http.StatusBadRequest, "Invalid state parameter.")
				return
			}

			// Exchange an authorization code for a token.
			token, err := authenticatorService.GetOath2Config().Exchange(ctx.Request.Context(), ctx.Query("code"))
			if err != nil {
				ctx.String(http.StatusUnauthorized, "Failed to exchange an authorization code for a token.")
				return
			}

			idToken, err := authenticatorService.VerifyIDToken(ctx.Request.Context(), token)
			if err != nil {
				ctx.String(http.StatusInternalServerError, "Failed to verify ID Token.")
				return
			}

			var profile map[string]interface{}
			if err := idToken.Claims(&profile); err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}

			session.Set("access_token", token.AccessToken)
			session.Set("profile", profile)
			if err := session.Save(); err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}

			// Redirect to logged in page.
			ctx.Redirect(http.StatusTemporaryRedirect, "/user")
		})

		// User
		r.GET("/user", IsAuthenticated, func(c *gin.Context) {
			session := sessions.Default(c)
			profile := session.Get("profile")

			c.HTML(http.StatusOK, "user.html", profile)
		})

		// Logout

		r.GET("/logout", func(c *gin.Context) {
			logoutUrl, err := url.Parse("https://" + authConf.GetDomain() + "/v2/logout")
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}

			scheme := "http"
			if c.Request.TLS != nil {
				scheme = "https"
			}

			returnTo, err := url.Parse(scheme + "://" + c.Request.Host)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}

			parameters := url.Values{}
			parameters.Add("returnTo", returnTo.String())
			parameters.Add("client_id", authConf.GetWebappClientId())
			logoutUrl.RawQuery = parameters.Encode()

			c.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
		})
	}
	var controllers []controller.Controller
	afterControllers := func(r *gin.Engine) {
		// Add gin engine configuration
	}
	coreAppServer := commonservers.NewCoreAppServerWithHooksImpl(
		serverConf,
		tlsConf,
		beforeControllers,
		controllers,
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
