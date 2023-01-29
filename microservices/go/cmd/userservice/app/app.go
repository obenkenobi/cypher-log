package app

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/background"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/servers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/lifecycle"
)

type App struct{}

func (a App) Start() {
	lifecycle.RunApp()

}

func NewApp(_ servers.GrpcServer, _ servers.AppServer, _ background.CronRunner) *App {
	return &App{}
}
