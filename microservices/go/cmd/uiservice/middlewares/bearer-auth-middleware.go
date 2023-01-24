package middlewares

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/security"
)

type BearerAuthMiddleware interface {
	PassBearerTokenFromSession() gin.HandlerFunc
}

type BearerAuthMiddlewareImpl struct {
}

func (b BearerAuthMiddlewareImpl) PassBearerTokenFromSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := sessions.Default(c).Get(security.AccessTokenSessionKey)
		c.Request.Header["Authorization"] = []string{fmt.Sprintf("Bearer %v", token)}
		c.Next()
	}
}

func NewBearerAuthMiddlewareImpl() *BearerAuthMiddlewareImpl {
	return &BearerAuthMiddlewareImpl{}
}
