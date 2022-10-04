package ginservices

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"sync"
)

type GinRouterProvider interface {
	AccessRouter(accessor func(router *gin.Engine))
}

type GinRouterProviderImpl struct {
	routerMu sync.Mutex
	router   *gin.Engine
}

func (g *GinRouterProviderImpl) AccessRouter(accessor func(router *gin.Engine)) {
	g.routerMu.Lock()
	defer g.routerMu.Unlock()
	accessor(g.router)
}

func NewGinEngineServiceImpl() *GinRouterProviderImpl {
	if environment.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	return &GinRouterProviderImpl{router: r}
}
