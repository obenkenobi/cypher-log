package servers

import (
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
) *AppServerImpl {
	coreAppServer := commonservers.NewCoreAppServerImpl(serverConf, tlsConf)
	return &AppServerImpl{CoreAppServer: coreAppServer}
}
