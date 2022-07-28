package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/services/go/pkg/logging"
	ginlogrus "github.com/toorop/gin-logrus"
)

func AddGlobalMiddleWares(r *gin.Engine) {
	log := logging.NewLogger()
	r.Use(ginlogrus.Logger(log), gin.Recovery())
}
