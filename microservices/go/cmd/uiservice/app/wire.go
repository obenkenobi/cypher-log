//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/controllers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/servers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf/authconf"
)

func InitializeApp() *App {
	wire.Build(
		conf.NewServerConfImpl,
		wire.Bind(new(conf.ServerConf), new(*conf.ServerConfImpl)),
		conf.NewTlsConfImpl,
		wire.Bind(new(conf.TLSConf), new(*conf.TlsConfImpl)),
		authconf.NewAuth0SecurityConfImpl,
		wire.Bind(new(authconf.Auth0SecurityConf), new(*authconf.Auth0RouteSecurityConfImpl)),
		services.NewAuthenticatorServiceImpl,
		wire.Bind(new(services.AuthenticatorService), new(*services.AuthenticatorServiceImpl)),
		middlewares.NewSessionMiddlewareImpl,
		wire.Bind(new(middlewares.SessionMiddleware), new(*middlewares.SessionMiddlewareImpl)),
		middlewares.NewBearerAuthMiddlewareImpl,
		wire.Bind(new(middlewares.BearerAuthMiddleware), new(*middlewares.BearerAuthMiddlewareImpl)),
		controllers.NewAuthControllerImpl,
		wire.Bind(new(controllers.AuthController), new(*controllers.AuthControllerImpl)),
		servers.NewAppServerImpl,
		wire.Bind(new(servers.AppServer), new(*servers.AppServerImpl)),
		NewApp)
	return &App{}
}
