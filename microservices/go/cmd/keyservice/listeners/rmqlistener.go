package listeners

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq/exchanges"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/rmqservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/taskrunner"
	"github.com/wagslane/go-rabbitmq"
)

type RmqListener interface {
	taskrunner.TaskRunner
}

type rmqListenerImpl struct {
	connector   rmqservices.RabbitConnector
	userService sharedservices.UserService
	ctx         context.Context
}

func (r rmqListenerImpl) ListenUserSave() {
	userCreateReceiver := rmq.NewReceiver(
		r.connector.GetConsumer(),
		"key_service_user_save",
		rmq.RoutingKeysDefault,
		"",
		exchanges.UserSaveExchange,
		rabbitmq.WithConsumeOptionsConcurrency(10),
		rabbitmq.WithConsumeOptionsQueueDurable,
		rabbitmq.WithConsumeOptionsQuorum,
	)
	userCreateReceiver.Listen(func(userDto userdtos.DistUserSaveDto) error {
		_, err := single.RetrieveValue(r.ctx, r.userService.SaveUser(r.ctx, userDto))
		return err
	}, true)
	logger.Log.Info("Listening for user saves")
}

func (r rmqListenerImpl) ListenUserDelete() {
	userCreateReceiver := rmq.NewReceiver(
		r.connector.GetConsumer(),
		"key_service_user_delete",
		rmq.RoutingKeysDefault,
		"",
		exchanges.UserDeleteExchange,
		rabbitmq.WithConsumeOptionsConcurrency(10),
		rabbitmq.WithConsumeOptionsQueueDurable,
		rabbitmq.WithConsumeOptionsQuorum,
	)
	userCreateReceiver.Listen(func(userDto userdtos.DistUserDeleteDto) error {
		_, err := single.RetrieveValue(r.ctx, r.userService.DeleteUser(r.ctx, userDto))
		return err
	}, true)
	logger.Log.Info("Listening for user deletions")
}

func (r rmqListenerImpl) Run() {
	forever := make(chan any)
	r.ListenUserSave()
	r.ListenUserDelete()
	<-forever
}

func NewRmqListener(connector rmqservices.RabbitConnector, userService sharedservices.UserService) RmqListener {
	return &rmqListenerImpl{ctx: context.Background(), connector: connector, userService: userService}
}
