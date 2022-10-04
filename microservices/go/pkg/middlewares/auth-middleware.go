package middlewares

import (
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/securityservices"
	"net/http"
)

// AuthorizerSettings provides settings on what criteria will be used to
// authorize access to an HTTP request.
type AuthorizerSettings struct {
	VerifyAnonymous        bool     // Verifies if the identity is anonymous
	VerifyIsUser           bool     // Verifies if the identity is a user
	VerifyIsSystemClient   bool     // Verifies if the identity is a system client
	AnyAuthoritiesToVerify []string // List of authorities of which at least one needs to be in the user's identity
	AllAuthoritiesToVerify []string // List of authorities of which all of them need to be in the user's identity
}

// AuthMiddleware provides Gin handler functions to authenticate and authorize a user.
type AuthMiddleware interface {
	// Authentication returns a handler to authenticate an HTTP request. If JWT
	// authentication is implemented, a JWT token will be verified.
	Authentication() gin.HandlerFunc

	// Authorization returns a handler that authorized a user using the provided AuthorizerSettings.
	// This handler should ONLY be called after the authentication handler is added to an HTTP router.
	Authorization(settings AuthorizerSettings) gin.HandlerFunc
}

type AuthMiddlewareImpl struct {
	jwtMiddleware        *jwtmiddleware.JWTMiddleware
	authorizationHandler gin.HandlerFunc
}

// NewAuthMiddlewareImpl creates an AuthMiddleware instance
func NewAuthMiddlewareImpl(
	externalOath2ValidateService securityservices.ExternalOath2ValidateService,
) *AuthMiddlewareImpl {

	jwtMiddleware := jwtmiddleware.New(externalOath2ValidateService.GetJwtValidator().ValidateToken)
	authorizationHandler := adapter.Wrap(jwtMiddleware.CheckJWT)
	return &AuthMiddlewareImpl{
		jwtMiddleware:        jwtMiddleware,
		authorizationHandler: authorizationHandler,
	}
}

func (a *AuthMiddlewareImpl) Authentication() gin.HandlerFunc {
	return a.authorizationHandler
}

func (a *AuthMiddlewareImpl) Authorization(settings AuthorizerSettings) gin.HandlerFunc {
	return func(c *gin.Context) {
		identity := security.GetIdentityFromGinContext(c)
		if (settings.VerifyAnonymous && !identity.IsAnonymous()) ||
			(settings.VerifyIsUser && identity.IsSystemClient()) ||
			(settings.VerifyIsSystemClient && identity.IsUser()) ||
			!identity.ContainsAnyAuthorities(settings.AnyAuthoritiesToVerify) ||
			!identity.ContainsAllAuthorities(settings.AllAuthoritiesToVerify) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
