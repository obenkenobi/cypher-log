package rmq

import amqp "github.com/rabbitmq/amqp091-go"

type ExchangeOptions[T any] struct {
	Name        string
	Kind        string
	Durable     bool
	AutoDeleted bool
	Internal    bool
	NoWait      bool
	Args        amqp.Table
}

type QueueOptions struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type ConsumeOptions struct {
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

type BindingOptions struct {
	RoutingKey string
	BindNoWait bool
	BindArgs   amqp.Table
}

type ReceiverOptions[T any] struct {
	ExchangeOpts       ExchangeOptions[T]
	QueueOptions       QueueOptions
	ConsumeOptions     ConsumeOptions
	BindingOptionsList []BindingOptions
	Uri                string
}

type SenderOptions[T any] struct {
	ExchangeOpts ExchangeOptions[T]
	RoutingKey   string
	Mandatory    bool
	Immediate    bool
	Uri          string
}
