package middlewares

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/security"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
)

type BearerAuthMiddleware interface {
	PassBearerTokenFromSession() gin.HandlerFunc
}

type BearerAuthMiddlewareImpl struct {
	accessTokenStoreService services.AccessTokenStoreService
}

func (b BearerAuthMiddlewareImpl) PassBearerTokenFromSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		userId := fmt.Sprintf("%v", session.Get(security.UserIdSessionKey))
		tokenId := fmt.Sprintf("%v", session.Get(security.TokenIdSessionKey))

		token, err := b.accessTokenStoreService.GetToken(c, userId, tokenId)
		if err != nil {
			logger.Log.WithError(err).Warn("Continuing the request but with an empty bearer token")
		}

		c.Request.Header["Authorization"] = []string{fmt.Sprintf("Bearer %v", token)}
		c.Next()
	}
}

func NewBearerAuthMiddlewareImpl(
	accessTokenStoreService services.AccessTokenStoreService,
) *BearerAuthMiddlewareImpl {
	return &BearerAuthMiddlewareImpl{accessTokenStoreService: accessTokenStoreService}
}
