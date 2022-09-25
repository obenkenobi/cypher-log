package rmq

import (
	"context"
	"encoding/json"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"time"
)

type Sender[T any] struct {
	_exchangeDeclaredRWLock sync.RWMutex
	_exchangeDeclared       bool
	senderOptions           SenderOptions
}

func BuildRabbitMQSender[T any](rmqPublishOpts SenderOptions) messaging.Sender[T] {
	return &Sender[T]{
		senderOptions:     rmqPublishOpts,
		_exchangeDeclared: false,
	}
}

func (r *Sender[T]) declareExchange(ch *amqp.Channel) error {
	var exchangeDeclared bool
	r._exchangeDeclaredRWLock.RLock()
	exchangeDeclared = r._exchangeDeclared
	r._exchangeDeclaredRWLock.RUnlock()
	if !exchangeDeclared {
		var err error = nil
		r._exchangeDeclaredRWLock.Lock()
		if !r._exchangeDeclared {
			err = ch.ExchangeDeclare(
				r.senderOptions.exchangeOpts.name,
				r.senderOptions.exchangeOpts.kind,
				r.senderOptions.exchangeOpts.durable,
				r.senderOptions.exchangeOpts.autoDeleted,
				r.senderOptions.exchangeOpts.internal,
				r.senderOptions.exchangeOpts.noWait,
				r.senderOptions.exchangeOpts.args,
			)
			r._exchangeDeclared = err == nil
		}
		r._exchangeDeclaredRWLock.Unlock()
		return err
	}
	return nil
}

func (r *Sender[T]) Send(msg T) single.Single[T] {
	return single.FromSupplier(func() (T, error) {
		conn, err := amqp.Dial(r.senderOptions.uri)
		if err != nil {
			return msg, err
		}
		defer conn.Close()
		ch, err := conn.Channel()
		if err != nil {
			return msg, err
		}
		err = r.declareExchange(ch)
		if err != nil {
			return msg, err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var amqpMsg amqp.Publishing
		{
			var msgMatcher interface{} = msg
			switch v := msgMatcher.(type) {
			case string:
				amqpMsg = amqp.Publishing{ContentType: ContentTypePlainText, Body: []byte(v)}
			default:
				body, err := json.Marshal(msg)
				if err != nil {
					return msg, err
				}
				amqpMsg = amqp.Publishing{ContentType: ContentTypeJson, Body: body}
			}
		}
		err = ch.PublishWithContext(ctx,
			r.senderOptions.exchangeOpts.name, // exchange
			r.senderOptions.routingKey,        // routing key
			r.senderOptions.mandatory,         // mandatory
			r.senderOptions.immediate,         // immediate
			amqpMsg)
		return msg, err
	})

}
