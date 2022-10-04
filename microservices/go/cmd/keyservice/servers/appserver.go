package servers

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/controllers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/commonservers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
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
	testController controllers.TestController,
) *AppServerImpl {
	coreAppServer := commonservers.NewCoreAppServerImpl(serverConf, tlsConf, testController)
	return &AppServerImpl{CoreAppServer: coreAppServer}
}
