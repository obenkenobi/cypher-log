package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/security"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/randutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/controller"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"strings"
)

type AuthController interface {
	controller.Controller
}

type AuthControllerImpl struct {
	authenticatorService    services.AuthenticatorService
	auth0SecurityConf       authconf.Auth0SecurityConf
	accessTokenStoreService services.AccessTokenStoreService
}

func (a AuthControllerImpl) AddRoutes(r *gin.Engine) {
	authGroup := r.Group("/auth")
	authGroup.GET("/login", func(c *gin.Context) {
		// Generate a random state
		state, err := randutils.GenerateRandom32BytesStr()
		if err != nil {
			a.sendInternalServerError(c, err)
			return
		}

		// Save the state inside the session.
		session := sessions.Default(c)
		session.Set(security.StateSessionKey, state)
		if err := session.Save(); err != nil {
			a.sendInternalServerError(c, err)
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, a.authenticatorService.GetOath2Config().AuthCodeURL(state,
			oauth2.SetAuthURLParam("audience", a.auth0SecurityConf.GetApiAudience())))
	})

	authGroup.GET("/callback", func(c *gin.Context) {
		session := sessions.Default(c)
		if c.Query(security.StateSessionKey) != session.Get(security.StateSessionKey) {
			c.String(http.StatusBadRequest, "Invalid state parameter.")
			return
		}

		// Exchange an authorization code for a token.
		token, err := a.authenticatorService.GetOath2Config().Exchange(c.Request.Context(), c.Query("code"))
		if err != nil {
			c.String(http.StatusUnauthorized, "Failed to exchange an authorization code for a token")
			return
		}

		idToken, err := a.authenticatorService.VerifyIDToken(c.Request.Context(), token)
		if err != nil {
			a.sendInternalServerError(c, err)
			return
		}

		randomUUID, err := uuid.NewRandom()
		if err != nil {
			a.sendInternalServerError(c, err)
			return
		}

		tokenId := strings.Join([]string{idToken.Subject, randomUUID.String()}, "/")

		session.Set(security.TokenIdSessionKey, tokenId)
		err = a.accessTokenStoreService.StoreToken(c, tokenId, token.AccessToken)
		if err != nil {
			a.sendInternalServerError(c, err)
			return
		}

		if err := session.Save(); err != nil {
			a.sendInternalServerError(c, err)
			return
		}

		// Redirect to logged in page.
		c.Redirect(http.StatusTemporaryRedirect, "/ui")
	})

	authGroup.GET("/logout", func(c *gin.Context) {
		logoutUrl, err := url.Parse("https://" + a.auth0SecurityConf.GetDomain() + "/v2/logout")
		if err != nil {
			a.sendInternalServerError(c, err)
			return
		}

		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}

		returnTo, err := url.Parse(scheme + "://" + c.Request.Host)
		if err != nil {
			logger.Log.WithError(err).Error()
			c.String(http.StatusInternalServerError, "Internal Server Error.")
			return
		}

		parameters := url.Values{}
		parameters.Add("returnTo", returnTo.String())
		parameters.Add("client_id", a.auth0SecurityConf.GetWebappClientId())
		logoutUrl.RawQuery = parameters.Encode()

		session := sessions.Default(c)
		session.Clear()
		if err := session.Save(); err != nil {
			a.sendInternalServerError(c, err)
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
	})
}

func (a AuthControllerImpl) sendInternalServerError(c *gin.Context, err error) {
	logger.Log.WithError(err).Error()
	c.String(http.StatusInternalServerError, "Internal Server Error.")
}

func NewAuthControllerImpl(
	authenticatorService services.AuthenticatorService,
	auth0SecurityConf authconf.Auth0SecurityConf,
	accessTokenStoreService services.AccessTokenStoreService,
) *AuthControllerImpl {
	return &AuthControllerImpl{
		accessTokenStoreService: accessTokenStoreService,
		authenticatorService:    authenticatorService,
		auth0SecurityConf:       auth0SecurityConf,
	}
}
