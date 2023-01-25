package middlewares

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/security"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/services"
)

type BearerAuthMiddleware interface {
	PassBearerTokenFromSession() gin.HandlerFunc
}

type BearerAuthMiddlewareImpl struct {
	accessTokenStoreService services.AccessTokenStoreService
}

func (b BearerAuthMiddlewareImpl) PassBearerTokenFromSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := sessions.Default(c).Get(security.AccessTokenSessionKey)
		c.Request.Header["Authorization"] = []string{fmt.Sprintf("Bearer %v", token)}
		c.Next()
	}
}

func NewBearerAuthMiddlewareImpl(
	accessTokenStoreService services.AccessTokenStoreService,
) *BearerAuthMiddlewareImpl {
	return &BearerAuthMiddlewareImpl{accessTokenStoreService: accessTokenStoreService}
}
