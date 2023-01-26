package rmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
	"github.com/wagslane/go-rabbitmq"
)

func readBody[T any](ctx context.Context, delivery rabbitmq.Delivery) (T, error) {
	isTypeString := false
	{
		var defaultBodyType T
		var matcher interface{} = defaultBodyType
		_, isTypeString = matcher.(string)
	}
	if delivery.ContentType == ContentTypePlainText && isTypeString {
		bodyStr := string(delivery.Body)
		var bodyAsInterface interface{} = bodyStr
		body, ok := bodyAsInterface.(T)
		var err error = nil
		if !ok {
			err = fmt.Errorf("failed to read consumed value as string")
			logger.Log.WithContext(ctx).Error(err)
		}
		return body, err
	}
	var body T
	err := json.Unmarshal(delivery.Body, &body)
	return body, err
}

type Receiver[T any] struct {
	consumer     rabbitmq.Consumer
	queue        string
	routingKeys  []string
	consumerName string
	exchange     Exchange[T]
	consumeOpts  []func(options *rabbitmq.ConsumeOptions)
}

func (r Receiver[T]) Listen(listener func(ctx context.Context, delivery messaging.Delivery[T]) messaging.ReceiverAction) {
	consumeOpts := append(r.exchange.GetConsumeOptions(), r.consumeOpts...)
	if utils.StringIsNotBlank(r.consumerName) {
		consumeOpts = append(consumeOpts, rabbitmq.WithConsumeOptionsConsumerName(r.consumerName))
	}
	err := r.consumer.StartConsuming(
		func(d rabbitmq.Delivery) rabbitmq.Action {
			ctx := context.Background()
			body, err := readBody[T](ctx, d)
			if err != nil {
				return rabbitmq.NackDiscard
			}
			action := listener(ctx, messaging.NewDelivery(body))
			switch action {
			case messaging.Commit:
				return rabbitmq.Ack
			case messaging.Discard:
				return rabbitmq.NackDiscard
			case messaging.Resend:
				return rabbitmq.NackRequeue
			default:
				return rabbitmq.Action(action)
			}
		},
		r.queue,
		r.routingKeys,
		consumeOpts...)
	if err != nil {
		logger.Log.Fatal(err)
	}
}

func NewReceiver[T any](
	consumer rabbitmq.Consumer,
	queue string,
	routingKeys []string,
	consumerName string,
	exchange Exchange[T],
	consumeOpts ...func(options *rabbitmq.ConsumeOptions)) messaging.Receiver[T] {
	return Receiver[T]{
		consumer:     consumer,
		queue:        queue,
		routingKeys:  routingKeys,
		consumerName: consumerName,
		exchange:     exchange,
		consumeOpts:  consumeOpts,
	}
}
