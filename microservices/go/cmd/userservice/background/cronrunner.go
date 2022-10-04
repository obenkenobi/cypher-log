package background

import (
	"context"
	"github.com/go-co-op/gocron"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/taskrunner"
	"time"
)

// CronRunner runs cron tasks in the background
type CronRunner interface {
	taskrunner.TaskRunner
}

type CronRunnerImpl struct {
	userService services.UserService
	ctx         context.Context
}

func (c CronRunnerImpl) Run() {
	s := gocron.NewScheduler(time.UTC)

	userChangeJob, err := s.Every(1).Second().Do(func() { c.userService.UsersChangeTask(c.ctx) })
	if err != nil {
		logger.Log.Fatal(err)
	}
	userChangeJob.SingletonMode()

	s.StartBlocking()
}

func NewCronRunnerImpl(userService services.UserService) *CronRunnerImpl {
	return &CronRunnerImpl{userService: userService, ctx: context.Background()}
}
