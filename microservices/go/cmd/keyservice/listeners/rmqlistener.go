package listeners

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	msg "github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq/exchanges"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/businessobjects/userbos"
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

func (r rmqListenerImpl) ListenUserChange() {
	userCreateReceiver := rmq.NewReceiver(
		r.connector.GetConsumer(),
		"key_service_user_change",
		rmq.RoutingKeysDefault,
		rmq.ConsumerDefault,
		exchanges.UserChangeExchange,
		rabbitmq.WithConsumeOptionsConcurrency(10),
		rabbitmq.WithConsumeOptionsQueueDurable,
		rabbitmq.WithConsumeOptionsQuorum,
	)
	userCreateReceiver.Listen(func(d msg.Delivery[userdtos.UserChangeEventDto]) msg.ReceiverAction {
		var userActionSrc single.Single[userbos.UserBo]
		switch d.Body().Action {
		case userdtos.UserSave:
			userActionSrc = r.userService.SaveUser(r.ctx, d.Body())
		case userdtos.UserDelete:
			userActionSrc = r.userService.DeleteUser(r.ctx, d.Body())
		default:
			return d.Discard()
		}
		if _, err := single.RetrieveValue(r.ctx, userActionSrc); err != nil {
			d.Resend()
		}
		return d.Commit()
	})
	logger.Log.Info("Listening for user changes")
}

func (r rmqListenerImpl) Run() {
	r.ListenUserChange()
	forever := make(chan any)
	<-forever
}

func NewRmqListener(connector rmqservices.RabbitConnector, userService sharedservices.UserService) RmqListener {
	return &rmqListenerImpl{ctx: context.Background(), connector: connector, userService: userService}
}
