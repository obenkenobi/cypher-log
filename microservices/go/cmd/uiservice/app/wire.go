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
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/ginservices"
)

func InitializeApp() *App {
	wire.Build(
		conf.NewExternalAppServerConfImpl,
		wire.Bind(new(conf.ExternalAppServerConf), new(*conf.ExternalAppServerConfImpl)),
		conf.NewSessionConfImpl,
		wire.Bind(new(conf.SessionConf), new(*conf.SessionConfImpl)),
		conf.NewStaticFilesConfImpl,
		wire.Bind(new(conf.StaticFilesConf), new(*conf.StaticFilesConfImpl)),
		conf.NewServerConfImpl,
		wire.Bind(new(conf.ServerConf), new(*conf.ServerConfImpl)),
		conf.NewTlsConfImpl,
		wire.Bind(new(conf.TLSConf), new(*conf.TlsConfImpl)),
		authconf.NewAuth0SecurityConfImpl,
		wire.Bind(new(authconf.Auth0SecurityConf), new(*authconf.Auth0RouteSecurityConfImpl)),
		sharedservices.NewErrorServiceImpl,
		wire.Bind(new(sharedservices.ErrorService), new(*sharedservices.ErrorServiceImpl)),
		ginservices.NewGinCtxServiceImpl,
		wire.Bind(new(ginservices.GinCtxService), new(*ginservices.GinCtxServiceImpl)),
		services.NewAuthenticatorServiceImpl,
		wire.Bind(new(services.AuthenticatorService), new(*services.AuthenticatorServiceImpl)),
		middlewares.NewSessionMiddlewareImpl,
		wire.Bind(new(middlewares.SessionMiddleware), new(*middlewares.SessionMiddlewareImpl)),
		middlewares.NewBearerAuthMiddlewareImpl,
		wire.Bind(new(middlewares.BearerAuthMiddleware), new(*middlewares.BearerAuthMiddlewareImpl)),
		middlewares.NewUiProviderMiddlewareImpl,
		wire.Bind(new(middlewares.UiProviderMiddleware), new(*middlewares.UiProviderMiddlewareImpl)),
		middlewares.NewUserKeyMiddlewareImpl,
		wire.Bind(new(middlewares.UserKeyMiddleware), new(*middlewares.UserKeyMiddlewareImpl)),
		controllers.NewAuthControllerImpl,
		wire.Bind(new(controllers.AuthController), new(*controllers.AuthControllerImpl)),
		controllers.NewCsrfControllerImpl,
		wire.Bind(new(controllers.CsrfController), new(*controllers.CsrfControllerImpl)),
		controllers.NewGatewayControllerImpl,
		wire.Bind(new(controllers.GatewayController), new(*controllers.GatewayControllerImpl)),
		servers.NewAppServerImpl,
		wire.Bind(new(servers.AppServer), new(*servers.AppServerImpl)),
		NewApp)
	return &App{}
}
