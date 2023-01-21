package servers

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/commonservers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/controller"
	"net/http"
)

type AppServer interface {
	commonservers.CoreAppServer
}

type AppServerImpl struct {
	commonservers.CoreAppServer
}

func NewAppServerImpl(
	serverConf conf.ServerConf,
	tlsConf conf.TLSConf,
) *AppServerImpl {
	beforeControllers := func(r *gin.Engine) {
		// Add gin engine configuration
		r.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "Hello World!!!")
		})
	}
	var controllers []controller.Controller
	afterControllers := func(r *gin.Engine) {
		// Add gin engine configuration
	}
	coreAppServer := commonservers.NewCoreAppServerWithHooksImpl(
		serverConf,
		tlsConf,
		beforeControllers,
		controllers,
		afterControllers,
	)
	return &AppServerImpl{CoreAppServer: coreAppServer}
}
