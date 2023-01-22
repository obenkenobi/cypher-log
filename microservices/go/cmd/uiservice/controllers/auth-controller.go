package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/security"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/randutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/controller"
	"net/http"
	"net/url"
)

type AuthController interface {
	controller.Controller
}

type AuthControllerImpl struct {
	authenticatorService services.AuthenticatorService
	auth0SecurityConf    authconf.Auth0SecurityConf
}

func (a AuthControllerImpl) AddRoutes(r *gin.Engine) {
	authGroup := r.Group("/auth")
	authGroup.GET("/login", func(ctx *gin.Context) {
		// Generate a random state
		state, err := randutils.GenerateRandom32Bytes()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Save the state inside the session.
		session := sessions.Default(ctx)
		session.Set(security.StateSessionKey, state)
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.Redirect(http.StatusTemporaryRedirect, a.authenticatorService.GetOath2Config().AuthCodeURL(state))
	})

	authGroup.GET("/callback", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		if ctx.Query(security.StateSessionKey) != session.Get(security.StateSessionKey) {
			ctx.String(http.StatusBadRequest, "Invalid state parameter.")
			return
		}

		// Exchange an authorization code for a token.
		token, err := a.authenticatorService.GetOath2Config().Exchange(ctx.Request.Context(), ctx.Query("code"))
		if err != nil {
			ctx.String(http.StatusUnauthorized, "Failed to exchange an authorization code for a token.")
			return
		}

		idToken, err := a.authenticatorService.VerifyIDToken(ctx.Request.Context(), token)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to verify ID Token.")
			return
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		session.Set(security.AccessTokenSessionKey, token.AccessToken)
		session.Set(security.ProfileSessionKey, profile)
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Redirect to logged in page.
		ctx.Redirect(http.StatusTemporaryRedirect, "/user")
	})

	authGroup.GET("/logout", func(c *gin.Context) {
		logoutUrl, err := url.Parse("https://" + a.auth0SecurityConf.GetDomain() + "/v2/logout")
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
		parameters.Add("client_id", a.auth0SecurityConf.GetWebappClientId())
		logoutUrl.RawQuery = parameters.Encode()

		c.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
	})
}

func NewAuthControllerImpl(
	authenticatorService services.AuthenticatorService,
	auth0SecurityConf authconf.Auth0SecurityConf,
) *AuthControllerImpl {
	return &AuthControllerImpl{
		authenticatorService: authenticatorService,
		auth0SecurityConf:    auth0SecurityConf,
	}
}
