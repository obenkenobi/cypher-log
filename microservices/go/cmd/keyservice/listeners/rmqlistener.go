package listeners

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	msg "github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq/exchanges"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/rmqservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/taskrunner"
	"github.com/wagslane/go-rabbitmq"
)

type RmqListener interface {
	taskrunner.TaskRunner
}

type rmqListenerImpl struct {
	connector              rmqservices.RabbitConnector
	userChangeEventService services.UserChangeEventService
	ctx                    context.Context
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
		resSrc := r.userChangeEventService.HandleUserChangeEventTransaction(r.ctx, d.Body())
		res, err := single.RetrieveValue(r.ctx, resSrc)
		if err != nil {
			return d.Resend()
		} else if res.Discarded {
			return d.Discard()
		} else {
			return d.Commit()
		}
	})
	logger.Log.Info("Listening for user changes")
}

func (r rmqListenerImpl) Run() {
	r.ListenUserChange()
	forever := make(chan any)
	<-forever
}

func NewRmqListener(
	connector rmqservices.RabbitConnector,
	userChangeEventService services.UserChangeEventService,
) RmqListener {
	return &rmqListenerImpl{
		ctx:                    context.Background(),
		connector:              connector,
		userChangeEventService: userChangeEventService,
	}
}
