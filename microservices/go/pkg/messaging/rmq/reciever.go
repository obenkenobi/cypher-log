package rmq

import (
	"encoding/json"
	"fmt"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
	"github.com/wagslane/go-rabbitmq"
)

func readBody[T any](delivery rabbitmq.Delivery) (T, error) {
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

func (r Receiver[T]) Listen(listener func(body T) error, resendIfErr bool) {
	consumeOpts := append(r.exchange.GetConsumeOptions(), r.consumeOpts...)
	if utils.StringIsNotBlank(r.consumerName) {
		consumeOpts = append(consumeOpts, rabbitmq.WithConsumeOptionsConsumerName(r.consumerName))
	}
	err := r.consumer.StartConsuming(
		func(d rabbitmq.Delivery) rabbitmq.Action {
			body, err := readBody[T](d)
			if err == nil {
				err = listener(body)
			}
			if err != nil {
				logger.Log.Error(err)
				if resendIfErr {
					return rabbitmq.NackRequeue
				}
				return rabbitmq.NackDiscard
			}
			logger.Log.Debugf("consumed: %v", string(d.Body))
			return rabbitmq.Ack
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
