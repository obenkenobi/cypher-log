package app

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/servers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/taskrunner"
)

type App struct {
	appServer servers.AppServer
}

func (a App) Start() {
	var taskRunners []taskrunner.TaskRunner

	if environment.ActivateAppServer() { // Add app server
		taskRunners = append(taskRunners, a.appServer)
	}

	taskrunner.RunAndWait(taskRunners...)
}

func NewApp(
	appServer servers.AppServer,
) *App {
	return &App{
		appServer: appServer,
	}
}
