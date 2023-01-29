package servers

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/controllers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/commonservers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/lifecycle"
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
	if !environment.ActivateAppServer() {
		// App server is deactivated, ran via the lifecycle package,
		// and is a root-child dependency so a nil is returned
		return nil
	}
	coreAppServer := commonservers.NewCoreAppServerImpl(serverConf, tlsConf, noteController)
	a := &AppServerImpl{CoreAppServer: coreAppServer}
	lifecycle.RegisterTaskRunner(a)
	return a
}
