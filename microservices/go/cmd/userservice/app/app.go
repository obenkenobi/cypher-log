package app

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/background"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/servers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/rmqservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/taskrunner"
)

type App struct {
	rmqPublisher rmqservices.RabbitMQPublisher
	grpcServer   servers.GrpcServer
	appServer    servers.AppServer
	cronRunner   background.CronRunner
}

func (a App) Start() {
	defer a.rmqPublisher.Close()
	// Add task dependencies
	var taskRunners []taskrunner.TaskRunner
	if environment.ActivateAppServer() { // Add app server
		taskRunners = append(taskRunners, a.appServer)
	}
	if environment.ActivateGrpcServer() { // Add GRPC server
		taskRunners = append(taskRunners, a.grpcServer)
	}
	if environment.ActivateCronRunner() { // Add Cron runner
		taskRunners = append(taskRunners, a.cronRunner)
	}
	// Run taskRunners
	taskrunner.RunAndWait(taskRunners...)

}

func NewApp(
	rmqPublisher rmqservices.RabbitMQPublisher,
	grpcServer servers.GrpcServer,
	appServer servers.AppServer,
	cronRunner background.CronRunner,
) *App {
	return &App{rmqPublisher: rmqPublisher, grpcServer: grpcServer, appServer: appServer, cronRunner: cronRunner}
}
