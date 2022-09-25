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
	senderOptions           SenderOptions[T]
}

func NewRabbitMQSender[T any](rmqPublishOpts SenderOptions[T]) messaging.Sender[T] {
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
				r.senderOptions.ExchangeOpts.Name,
				r.senderOptions.ExchangeOpts.Kind,
				r.senderOptions.ExchangeOpts.Durable,
				r.senderOptions.ExchangeOpts.AutoDeleted,
				r.senderOptions.ExchangeOpts.Internal,
				r.senderOptions.ExchangeOpts.NoWait,
				r.senderOptions.ExchangeOpts.Args,
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
		conn, err := amqp.Dial(r.senderOptions.Uri)
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
			r.senderOptions.ExchangeOpts.Name, // exchange
			r.senderOptions.RoutingKey,        // routing key
			r.senderOptions.Mandatory,         // mandatory
			r.senderOptions.Immediate,         // immediate
			amqpMsg)
		return msg, err
	})

}
