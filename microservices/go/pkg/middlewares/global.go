package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	ginLog "github.com/toorop/gin-logrus"
)

func AddGlobalMiddleWares(r *gin.Engine) {
	r.Use(ginLog.Logger(logger.Log), gin.Recovery())
}
