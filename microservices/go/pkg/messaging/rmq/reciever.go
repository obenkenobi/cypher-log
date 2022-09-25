package rmq

import (
	"encoding/json"
	"github.com/barweiss/go-tuple"
	"github.com/joamaki/goreactive/stream"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type delivered[T any] struct {
	body      T
	delivery  amqp.Delivery
	isAutoAck bool
}

func (d delivered[T]) GetBody() T {
	return d.body
}

func (d delivered[T]) Ack(multiple bool) error {
	if d.IsAutoAck() {
		return nil
	}
	return d.delivery.Ack(multiple)
}

func (d delivered[T]) Reject(sendBack bool) error {
	if d.IsAutoAck() {
		return nil
	}
	return d.delivery.Reject(sendBack)
}

func (d delivered[T]) IsAutoAck() bool { return d.isAutoAck }

func createDelivered[T any](delivery amqp.Delivery, isAutoAck bool) (messaging.Delivered[T], error) {
	isTypeString := false
	{
		var defaultBodyType T
		var matcher interface{} = defaultBodyType
		_, isTypeString = matcher.(string)

	}
	if delivery.ContentType == ContentTypePlainText && isTypeString {
		bodyStr := string(delivery.Body)
		var bodyAsInterface interface{} = bodyStr
		return delivered[T]{body: bodyAsInterface.(T), delivery: delivery, isAutoAck: isAutoAck}, nil
	}
	var body T
	if err := json.Unmarshal(delivery.Body, &body); err != nil {
		return nil, err
	}
	return delivered[T]{body: body, delivery: delivery, isAutoAck: isAutoAck}, nil

}

type Receiver[T any] struct {
	err        error
	conn       *amqp.Connection
	ch         *amqp.Channel
	autoAck    bool
	messageSrc stream.Observable[amqp.Delivery]
}

func (r Receiver[T]) IsAutoAck() bool { return r.autoAck }

func (r Receiver[T]) MessageStream() stream.Observable[T] {
	deliveredObs := stream.Map(r.messageSrc, func(d amqp.Delivery) tuple.T2[messaging.Delivered[T], error] {
		msgDelivered, err := createDelivered[T](d, r.IsAutoAck())
		return tuple.New2(msgDelivered, err)
	})
	successMsgObs := stream.Filter(deliveredObs, func(pair tuple.T2[messaging.Delivered[T], error]) bool {
		err := pair.V2
		isSuccess := err == nil
		if isSuccess {
			return true
		} else {
			log.Error(err)
			return false
		}
	})
	return stream.Map(successMsgObs, func(pair tuple.T2[messaging.Delivered[T], error]) messaging.Delivered[T] {
		return pair.V1
	})
}

func (r Receiver[T]) Close() {
	defer func(conn *amqp.Connection) {
		if conn != nil {
			return
		}
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}(r.conn)
	defer func(ch *amqp.Channel) {
		if r.ch != nil {
			return
		}
		err := ch.Close()
		if err != nil {
			log.Error(err)
		}
	}(r.ch)
}

// StartReceiver connects to a rabbitMQ instance and returns a receiver that can be listened to.
func StartReceiver[T any](receiverOpts ReceiverOptions[T]) messaging.Receiver[T] {
	var receiver = Receiver[T]{}
	conn, err := amqp.Dial(receiverOpts.Uri)
	receiver.conn = conn
	receiver.err = err
	if err != nil {
		return receiver
	}

	ch, err := conn.Channel()
	receiver.ch = ch
	receiver.err = err
	if err != nil {
		return receiver
	}

	err = ch.ExchangeDeclare(
		receiverOpts.ExchangeOpts.Name,
		receiverOpts.ExchangeOpts.Kind,
		receiverOpts.ExchangeOpts.Durable,
		receiverOpts.ExchangeOpts.AutoDeleted,
		receiverOpts.ExchangeOpts.Internal,
		receiverOpts.ExchangeOpts.NoWait,
		receiverOpts.ExchangeOpts.Args,
	)
	receiver.err = err
	if err != nil {
		return receiver
	}

	q, err := ch.QueueDeclare(
		receiverOpts.QueueOptions.Name,       // name
		receiverOpts.QueueOptions.Durable,    // durable
		receiverOpts.QueueOptions.AutoDelete, // delete when unused
		receiverOpts.QueueOptions.Exclusive,  // exclusive
		receiverOpts.QueueOptions.NoWait,     // no-wait
		receiverOpts.QueueOptions.Args,       // arguments
	)
	receiver.err = err
	if err != nil {
		return receiver
	}
	for _, bindingOpt := range receiverOpts.BindingOptionsList {
		err = ch.QueueBind(
			q.Name,                         // queue name
			bindingOpt.RoutingKey,          // routing key
			receiverOpts.ExchangeOpts.Name, // exchange
			bindingOpt.BindNoWait,
			bindingOpt.BindArgs)
		receiver.err = err
		if err != nil {
			return receiver
		}
	}
	messages, err := ch.Consume(
		q.Name,                                // queue
		receiverOpts.ConsumeOptions.Consumer,  // consumer
		receiverOpts.ConsumeOptions.AutoAck,   // auto ack
		receiverOpts.ConsumeOptions.Exclusive, // exclusive
		receiverOpts.ConsumeOptions.NoLocal,   // no local
		receiverOpts.ConsumeOptions.NoWait,    // no wait
		receiverOpts.ConsumeOptions.Args,      // args
	)
	receiver.err = err
	receiver.messageSrc = stream.FromChannel(messages)
	return receiver
}
