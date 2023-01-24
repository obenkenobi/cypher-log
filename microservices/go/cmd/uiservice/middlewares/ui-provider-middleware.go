package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
	"net/http"
)

type UiProviderMiddleware interface {
	ProvideUI(r *gin.Engine)
}

type UiProviderMiddlewareImpl struct {
	staticFilesConf conf.StaticFilesConf
}

func (u UiProviderMiddlewareImpl) ProvideUI(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/ui")
	})

	if environment.IsDevelopment() {
		r.GET("/ui/*proxyPath", func(c *gin.Context) {
			// Redirect to gatsby server in dev settings
			c.Redirect(http.StatusTemporaryRedirect, "https://localhost:8000/"+c.Param("proxyPath"))
		})
	} else {
		staticFilesPath := u.staticFilesConf.GetStaticFilesPath()
		if utils.StringIsBlank(staticFilesPath) {
			staticFilesPath = "cmd/uiservice/ClientApp/public"
		}

		r.Static("ui/", staticFilesPath)
	}
}

func NewUiProviderMiddlewareImpl(staticFilesConf conf.StaticFilesConf) *UiProviderMiddlewareImpl {
	return &UiProviderMiddlewareImpl{staticFilesConf: staticFilesConf}
}
