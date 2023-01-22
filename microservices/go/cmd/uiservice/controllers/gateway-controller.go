package controllers

import (
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/controller"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type GatewayController interface {
	controller.Controller
}

type GatewayControllerImpl struct {
	externalAppServerConf conf.ExternalAppServerConf
	tlsConf               conf.TLSConf
}

func (g GatewayControllerImpl) AddRoutes(r *gin.Engine) {
	apiGroup := r.Group("/api")

	apiGroup.Any("/userservice/*proxyPath",
		g.proxyHandler("proxyPath", g.externalAppServerConf.GetUserServiceAddress()))

	apiGroup.Any("/keyservice/*proxyPath",
		g.proxyHandler("proxyPath", g.externalAppServerConf.GetKeyServiceAddress()))

	apiGroup.Any("/noteservice/*proxyPath",
		g.proxyHandler("proxyPath", g.externalAppServerConf.GetNoteServiceAddress()))
}

func (g GatewayControllerImpl) proxyHandler(srcUrlParam string, destAddr string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, err := url.Parse(destAddr)
		if err != nil {
			panic(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)

		roundTripper, err := g.httpRoundTripper()
		if err != nil {
			logger.Log.WithError(err).Error()
			c.Status(http.StatusBadRequest)
			return
		}
		proxy.Transport = roundTripper

		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = c.Param(srcUrlParam)
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func (g GatewayControllerImpl) httpRoundTripper() (http.RoundTripper, error) {
	if !environment.IsProduction() {
		logger.Log.Info("Non production environments can skip TLS verification")
		return &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}, nil
	}
	logger.Log.Info("Must require a secure TLS connection")
	return http.DefaultTransport, nil
}
func NewGatewayControllerImpl(
	externalAppServerConf conf.ExternalAppServerConf,
	tlsConf conf.TLSConf,
) *GatewayControllerImpl {
	return &GatewayControllerImpl{externalAppServerConf: externalAppServerConf, tlsConf: tlsConf}
}
