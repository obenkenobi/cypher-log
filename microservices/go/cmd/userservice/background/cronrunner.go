package background

import (
	"context"
	"github.com/go-co-op/gocron"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/lifecycle"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"time"
)

// CronRunner runs cron tasks in the background
type CronRunner interface {
	lifecycle.TaskRunner
}

type CronRunnerImpl struct {
	userService services.UserService
}

func (c CronRunnerImpl) Run() {
	s := gocron.NewScheduler(time.UTC)

	userChangeJob, err := s.Every(1).Second().Do(func() {
		ctx := context.Background()
		c.userService.UsersChangeTask(ctx)
	})
	if err != nil {
		logger.Log.Fatal(err)
	}
	userChangeJob.SingletonMode()

	s.StartBlocking()
}

func NewCronRunnerImpl(userService services.UserService) *CronRunnerImpl {
	if !environment.ActivateCronRunner() {
		// Task runner is deactivated, ran via the lifecycle package,
		// and is a root-child dependency so a nil is returned
		return nil
	}
	c := &CronRunnerImpl{userService: userService}
	lifecycle.RegisterTaskRunner(c)
	return c
}
