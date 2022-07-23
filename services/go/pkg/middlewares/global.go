package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/services/go/pkg/logging"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

func AddGlobalMiddleWares(r *gin.Engine) {
	log := logrus.New()
	logging.ConfigTextLoggingWithLogger(log)
	r.Use(ginlogrus.Logger(log), gin.Recovery())
}
