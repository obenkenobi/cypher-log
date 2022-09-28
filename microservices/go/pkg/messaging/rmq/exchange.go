package rmq

import (
	"github.com/wagslane/go-rabbitmq"
)

type Exchange[T any] struct {
	Name        string
	Kind        string
	Durable     bool
	AutoDeleted bool
	Internal    bool
	NoWait      bool
}

func (e Exchange[T]) GetConsumeOptions() []func(*rabbitmq.ConsumeOptions) {
	consumeOpts := []func(*rabbitmq.ConsumeOptions){
		rabbitmq.WithConsumeOptionsBindingExchangeName(e.Name),
		rabbitmq.WithConsumeOptionsBindingExchangeKind(e.Kind),
	}
	if e.Durable {
		consumeOpts = append(consumeOpts, rabbitmq.WithConsumeOptionsBindingExchangeDurable)
	}
	if e.AutoDeleted {
		consumeOpts = append(consumeOpts, rabbitmq.WithConsumeOptionsBindingExchangeAutoDelete)
	}
	if e.Internal {
		consumeOpts = append(consumeOpts, rabbitmq.WithConsumeOptionsBindingExchangeInternal)
	}
	if e.NoWait {
		consumeOpts = append(consumeOpts, rabbitmq.WithConsumeOptionsBindingExchangeNoWait)
	}
	return consumeOpts
}

func (e Exchange[T]) GetPublishOptions() []func(options *rabbitmq.PublishOptions) {
	return []func(*rabbitmq.PublishOptions){rabbitmq.WithPublishOptionsExchange(e.Name)}
}
