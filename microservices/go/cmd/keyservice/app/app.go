package app

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/listeners"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/servers"
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

	// Add task dependencies
	var taskRunners []taskrunner.TaskRunner
	if environment.ActivateAppServer() { // Add app server
		taskRunners = append(taskRunners, a.appServer)
	}
	if environment.ActivateRabbitMqListener() {
		taskRunners = append(taskRunners, a.rmqListener)
	}
	// Run tasks
	taskrunner.RunAndWait(taskRunners...)
}

func NewApp(
	rabbitConsumer rmqservices.RabbitMQConsumer,
	appServer servers.AppServer,
	rmqListener listeners.RmqListener,
) *App {
	return &App{rabbitConsumer: rabbitConsumer, appServer: appServer, rmqListener: rmqListener}
}
