package listeners

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	msg "github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq/exchanges"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/rmqservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/taskrunner"
	"github.com/wagslane/go-rabbitmq"
)

type RmqListener interface {
	taskrunner.TaskRunner
}

type RmqListenerImpl struct {
	consumer               rmqservices.RabbitMQConsumer
	userChangeEventService services.UserChangeEventService
	ctx                    context.Context
}

func (r RmqListenerImpl) ListenUserChange() {
	userCreateReceiver := rmq.NewReceiver(
		r.consumer.GetConsumer(),
		"note_service_user_change",
		rmq.RoutingKeysDefault,
		rmq.ConsumerDefault,
		exchanges.UserChangeExchange,
		rabbitmq.WithConsumeOptionsConcurrency(10),
		rabbitmq.WithConsumeOptionsQueueDurable,
		rabbitmq.WithConsumeOptionsQuorum,
	)
	userCreateReceiver.Listen(func(ctx context.Context, d msg.Delivery[userdtos.UserChangeEventDto]) msg.ReceiverAction {
		res, err := r.userChangeEventService.HandleUserChangeEventTxn(ctx, d.Body())
		if err != nil {
			return d.Resend()
		} else if res.Discarded {
			return d.Discard()
		} else {
			return d.Commit()
		}
	})
	logger.Log.WithContext(r.ctx).Info("Listening for user changes")
}

func (r RmqListenerImpl) Run() {
	r.ListenUserChange()
	forever := make(chan any)
	<-forever
}

func NewRmqListenerImpl(
	consumer rmqservices.RabbitMQConsumer,
	userChangeEventService services.UserChangeEventService,
) *RmqListenerImpl {
	return &RmqListenerImpl{
		ctx:                    context.Background(),
		consumer:               consumer,
		userChangeEventService: userChangeEventService,
	}
}
