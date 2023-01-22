package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"net/http"
)

type UiProviderMiddleware interface {
	ProvideUI(r *gin.Engine)
}

type UiProviderMiddlewareImpl struct {
}

func (u UiProviderMiddlewareImpl) ProvideUI(r *gin.Engine) {
	//Todo: configure static path via env variable

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/ui")
	})

	if environment.IsDevelopment() {
		r.LoadHTMLGlob("cmd/uiservice/resources/web/template/*")

		r.GET("/ui", func(c *gin.Context) {
			c.HTML(http.StatusOK, "home.html", struct{}{})
		})
	} else {
		r.Static("ui/", "cmd/uiservice/ClientApp/public")
	}
}

func NewUiProviderMiddlewareImpl() *UiProviderMiddlewareImpl {
	return &UiProviderMiddlewareImpl{}
}
