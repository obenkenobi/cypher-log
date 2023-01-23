package controllers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/commondtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/controller"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

type GatewayController interface {
	controller.Controller
}

type GatewayControllerImpl struct {
	externalAppServerConf conf.ExternalAppServerConf
	bearerAuthMiddleware  middlewares.BearerAuthMiddleware
	tlsConf               conf.TLSConf
}

func (g GatewayControllerImpl) AddRoutes(r *gin.Engine) {
	apiGroup := r.Group("/api", g.bearerAuthMiddleware.PassBearerTokenFromSession())

	apiGroup.Any("/userservice/*proxyPath",
		g.proxyHandler("proxyPath", g.externalAppServerConf.GetUserServiceAddress()))

	apiGroup.Any("/keyservice/*proxyPath",
		g.proxyHandlerWithModifyResp("proxyPath", g.externalAppServerConf.GetKeyServiceAddress(),
			func(res *http.Response, destPath string, c *gin.Context) error {
				// Only apply to paths starting with "v1/userKey/newSession"
				if !strings.HasPrefix(destPath, "v1/userKey/newSession") {
					return nil
				}
				// original bytes to session dto
				originalBytes, err := ioutil.ReadAll(res.Body) //Read html
				if err != nil {
					return err
				}
				if err = res.Body.Close(); err != nil {
					return err
				}
				sessionDto := commondtos.UKeySessionDto{}
				if err = json.Unmarshal(originalBytes, &sessionDto); err != nil {
					return err
				}

				// write session sto to session
				security.WriteUKeySessionDtoToSession(c, sessionDto)

				// replace body with a success dto
				newBytes, err := json.Marshal(commondtos.NewSuccessTrue())
				if err != nil {
					return err
				}
				body := ioutil.NopCloser(bytes.NewReader(newBytes))
				res.Body = body
				res.ContentLength = int64(len(newBytes))
				res.Header.Set("Content-Length", strconv.Itoa(len(newBytes)))
				return nil
			},
		))

	apiGroup.Any("/noteservice/*proxyPath",
		g.proxyHandler("proxyPath", g.externalAppServerConf.GetNoteServiceAddress()))
}

func (g GatewayControllerImpl) proxyHandler(destPathUrlParam string, destAddr string) gin.HandlerFunc {
	return g.proxyHandlerWithModifyResp(destPathUrlParam, destAddr, nil)
}

func (g GatewayControllerImpl) proxyHandlerWithModifyResp(
	destPathUrlParam string,
	destAddr string,
	modifyRes func(res *http.Response, destPath string, c *gin.Context) error,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, err := url.Parse(destAddr)
		if err != nil {
			panic(err)
		}

		destPath := c.Param(destPathUrlParam)

		proxy := httputil.NewSingleHostReverseProxy(remote)

		if modifyRes != nil {
			proxy.ModifyResponse = func(res *http.Response) error {
				return modifyRes(res, destPath, c)
			}
		}

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
			req.URL.Path = destPath
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
	bearerAuthMiddleware middlewares.BearerAuthMiddleware,
	tlsConf conf.TLSConf,
) *GatewayControllerImpl {
	return &GatewayControllerImpl{
		bearerAuthMiddleware:  bearerAuthMiddleware,
		externalAppServerConf: externalAppServerConf,
		tlsConf:               tlsConf,
	}
}
