package app

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/servers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/lifecycle"
)

type App struct{}

func (a App) Start() {
	lifecycle.RunApp()
}

func NewApp(_ servers.AppServer) *App {
	return &App{}
}
