package servers

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/controllers"
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
	noteController controllers.NoteController,
) *AppServerImpl {
	coreAppServer := commonservers.NewCoreAppServerImpl(serverConf, tlsConf, noteController)
	return &AppServerImpl{CoreAppServer: coreAppServer}
}
