package middlewares

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/security"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
)

type BearerAuthMiddleware interface {
	PassBearerTokenFromSession() gin.HandlerFunc
}

type BearerAuthMiddlewareImpl struct {
	accessTokenStoreService services.AccessTokenStoreService
}

func (b BearerAuthMiddlewareImpl) PassBearerTokenFromSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Log.WithContext(c).Info("Starting bearer token middleware")
		tokenId := utils.AnyToString(sessions.Default(c).Get(security.TokenIdSessionKey))

		token, err := b.accessTokenStoreService.GetToken(c, tokenId)
		if err != nil {
			logger.Log.WithContext(c).WithError(err).Warn("Continuing the request but with an empty bearer token")
		}
		if utils.StringIsBlank(token) {
			token = "_"
		}

		c.Request.Header["Authorization"] = []string{fmt.Sprintf("Bearer %v", token)}
		logger.Log.WithContext(c).Info("Ending bearer token middleware")
		c.Next()
	}
}

func NewBearerAuthMiddlewareImpl(
	accessTokenStoreService services.AccessTokenStoreService,
) *BearerAuthMiddlewareImpl {
	return &BearerAuthMiddlewareImpl{accessTokenStoreService: accessTokenStoreService}
}
