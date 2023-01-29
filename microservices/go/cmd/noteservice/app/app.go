package app

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/listeners"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/servers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/lifecycle"
)

type App struct{}

func (a App) Start() {
	lifecycle.RunApp()
}

func NewApp(_ servers.AppServer, _ listeners.KafkaListener) *App {
	return &App{}
}
