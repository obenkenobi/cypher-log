package app

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/listeners"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/servers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/rmqservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/taskrunner"
)

type App struct {
	rabbitConsumer rmqservices.RabbitMQConsumer
	rmqListener    listeners.RmqListener
	appServer      servers.AppServer
}

func (a App) Start() {
	defer a.rabbitConsumer.Close()

	var taskRunners []taskrunner.TaskRunner
	if environment.ActivateAppServer() { // Add app server
		taskRunners = append(taskRunners, a.appServer)
	}
	if environment.ActivateRabbitMqListener() {
		taskRunners = append(taskRunners, a.rmqListener)
	}
	taskrunner.RunAndWait(taskRunners...)
}

func NewApp(
	appServer servers.AppServer,
	rmqListener listeners.RmqListener,
	rabbitConsumer rmqservices.RabbitMQConsumer,
) *App {
	return &App{
		rmqListener:    rmqListener,
		appServer:      appServer,
		rabbitConsumer: rabbitConsumer,
	}
}
