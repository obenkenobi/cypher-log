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

func CreateReceiver[T any](receiverOpts ReceiverOptions) messaging.Receiver[T] {
	var receiver = Receiver[T]{}
	conn, err := amqp.Dial(receiverOpts.uri)
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
		receiverOpts.exchangeOpts.name,
		receiverOpts.exchangeOpts.kind,
		receiverOpts.exchangeOpts.durable,
		receiverOpts.exchangeOpts.autoDeleted,
		receiverOpts.exchangeOpts.internal,
		receiverOpts.exchangeOpts.noWait,
		receiverOpts.exchangeOpts.args,
	)
	receiver.err = err
	if err != nil {
		return receiver
	}

	q, err := ch.QueueDeclare(
		receiverOpts.queueOptions.name,       // name
		receiverOpts.queueOptions.durable,    // durable
		receiverOpts.queueOptions.autoDelete, // delete when unused
		receiverOpts.queueOptions.exclusive,  // exclusive
		receiverOpts.queueOptions.noWait,     // no-wait
		receiverOpts.queueOptions.args,       // arguments
	)
	receiver.err = err
	if err != nil {
		return receiver
	}
	for _, routingKey := range receiverOpts.routingKeys {
		err = ch.QueueBind(
			q.Name,                         // queue name
			routingKey,                     // routing key
			receiverOpts.exchangeOpts.name, // exchange
			receiverOpts.bindNoWait,
			receiverOpts.bindArgs)
		receiver.err = err
		if err != nil {
			return receiver
		}
	}
	messages, err := ch.Consume(
		q.Name,                                // queue
		receiverOpts.consumeOptions.consumer,  // consumer
		receiverOpts.consumeOptions.autoAck,   // auto ack
		receiverOpts.consumeOptions.exclusive, // exclusive
		receiverOpts.consumeOptions.noLocal,   // no local
		receiverOpts.consumeOptions.noWait,    // no wait
		receiverOpts.consumeOptions.args,      // args
	)
	receiver.err = err
	receiver.messageSrc = stream.FromChannel(messages)
	return receiver
}
