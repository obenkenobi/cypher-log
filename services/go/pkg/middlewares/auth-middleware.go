package middlewares

import (
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/services/go/pkg/security"
	"net/http"
	"time"
)

type AuthorizerMiddlewareParams struct {
	RequireAnyAuthorities []string // List of authorities of which at least one needs to be in the user's identity
	RequireAllAuthorities []string // List of authorities of which all of them need to be in the user's identity
}

type AuthMiddleware interface {
	Authentication() gin.HandlerFunc
	Authorization(params AuthorizerMiddlewareParams) gin.HandlerFunc
}

type AuthMiddlewareImpl struct {
	provider             *jwks.CachingProvider
	jwtValidator         *validator.Validator
	jwtMiddleware        *jwtmiddleware.JWTMiddleware
	authorizationHandler gin.HandlerFunc
}

func BuildAuthMiddleware(auth0Conf conf.Auth0Conf) AuthMiddleware {
	issuerURL := auth0Conf.GetIssuerUrl()
	audience := auth0Conf.GetAudience()

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)
	jwtValidator, _ := validator.New(provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{audience},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &security.Auth0CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)

	jwtMiddleware := jwtmiddleware.New(jwtValidator.ValidateToken)
	authorizationHandler := adapter.Wrap(jwtMiddleware.CheckJWT)
	return &AuthMiddlewareImpl{
		provider:             provider,
		jwtValidator:         jwtValidator,
		jwtMiddleware:        jwtMiddleware,
		authorizationHandler: authorizationHandler,
	}
}

func (a *AuthMiddlewareImpl) Authentication() gin.HandlerFunc {
	return a.authorizationHandler
}

func (a *AuthMiddlewareImpl) Authorization(params AuthorizerMiddlewareParams) gin.HandlerFunc {
	return func(c *gin.Context) {
		identityHolder := security.NewIdentityHolderFromContext(c)
		if !identityHolder.ValidateAuthoritiesAny(params.RequireAnyAuthorities) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		} else if !identityHolder.ValidateAuthoritiesAll(params.RequireAllAuthorities) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
