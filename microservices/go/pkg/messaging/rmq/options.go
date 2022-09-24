package rmq

import amqp "github.com/rabbitmq/amqp091-go"

type ExchangeOptions struct {
	name        string
	kind        string
	durable     bool
	autoDeleted bool
	internal    bool
	noWait      bool
	args        amqp.Table
}

type QueueOptions struct {
	name       string
	durable    bool
	autoDelete bool
	exclusive  bool
	noWait     bool
	args       amqp.Table
}

type ConsumeOptions struct {
	consumer  string
	autoAck   bool
	exclusive bool
	noLocal   bool
	noWait    bool
	args      amqp.Table
}

type ReceiverOptions struct {
	exchangeOpts   ExchangeOptions
	queueOptions   QueueOptions
	consumeOptions ConsumeOptions
	routingKeys    []string
	bindNoWait     bool
	bindArgs       amqp.Table
	uri            string
}

type SenderOptions struct {
	exchangeOpts ExchangeOptions
	routingKey   string
	mandatory    bool
	immediate    bool
	uri          string
}
