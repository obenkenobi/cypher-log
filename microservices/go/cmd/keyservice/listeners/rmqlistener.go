package listeners

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq/exchanges"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq/rmqservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/taskrunner"
	"github.com/wagslane/go-rabbitmq"
)

type RmqListener interface {
	taskrunner.TaskRunner
}

type rmqListenerImpl struct {
	connector rmqservices.RabbitConnector
}

func (r rmqListenerImpl) ListenUserCreate() {
	userCreateReceiver := rmq.NewReceiver(
		r.connector.GetConsumer(),
		"key_service_user_create",
		rmq.RoutingKeysDefault,
		"",
		exchanges.UserCreateExchange,
		rabbitmq.WithConsumeOptionsConcurrency(10),
		rabbitmq.WithConsumeOptionsQueueDurable,
		rabbitmq.WithConsumeOptionsQuorum,
	)
	userCreateReceiver.Listen(func(userDto userdtos.UserDto) error {
		logger.Log.Info("creating user", userDto)
		return nil
	}, true)
	logger.Log.Info("Listening for user creation")
}

func (r rmqListenerImpl) ListenUserUpdate() {
	userCreateReceiver := rmq.NewReceiver(
		r.connector.GetConsumer(),
		"key_service_user_update",
		rmq.RoutingKeysDefault,
		"",
		exchanges.UserUpdateExchange,
		rabbitmq.WithConsumeOptionsConcurrency(10),
		rabbitmq.WithConsumeOptionsQueueDurable,
		rabbitmq.WithConsumeOptionsQuorum,
	)
	userCreateReceiver.Listen(func(userDto userdtos.UserDto) error {
		logger.Log.Info("updating user", userDto)
		return nil
	}, true)
	logger.Log.Info("Listening for user updates")
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
	userCreateReceiver.Listen(func(userDto userdtos.UserDto) error {
		logger.Log.Info("deleting user", userDto)
		return nil
	}, true)
	logger.Log.Info("Listening for user deletions")
}

func (r rmqListenerImpl) Run() {
	forever := make(chan any)
	go r.ListenUserCreate()
	go r.ListenUserUpdate()
	go r.ListenUserDelete()
	<-forever
}

func NewRmqListener(connector rmqservices.RabbitConnector) RmqListener {
	return &rmqListenerImpl{connector: connector}
}
