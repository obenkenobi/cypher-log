package rmqservices

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/wagslane/go-rabbitmq"
)

type RabbitMQPublisher interface {
	GetPublisher() *rabbitmq.Publisher
	Close()
}

type RabbitMQPublisherImpl struct {
	publisher *rabbitmq.Publisher
}

func (r RabbitMQPublisherImpl) GetPublisher() *rabbitmq.Publisher { return r.publisher }

func (r *RabbitMQPublisherImpl) Close() {
	if r != nil {
		defer func(publisher *rabbitmq.Publisher) {
			if publisher == nil {
				return
			}
			err := publisher.Close()
			if err != nil {
				logger.Log.Error(err)
			}
		}(r.publisher)
	}
}

func NewRabbitPublisherImpl(rabbitMQConf conf.RabbitMQConf) *RabbitMQPublisherImpl {
	publisher, err := rabbitmq.NewPublisher(
		rabbitMQConf.GetURI(),
		rabbitmq.Config{},
		// can pass nothing for no logging
		rabbitmq.WithPublisherOptionsLogger(logger.Log),
	)
	if err != nil {
		logger.Log.Fatal(err)
	}
	return &RabbitMQPublisherImpl{publisher: publisher}
}

type RabbitMQConsumer interface {
	GetConsumer() rabbitmq.Consumer
	Close()
}

type RabbitMQConsumerImpl struct {
	consumer rabbitmq.Consumer
}

func (r RabbitMQConsumerImpl) GetConsumer() rabbitmq.Consumer { return r.consumer }

func (r *RabbitMQConsumerImpl) Close() {
	if r != nil {
		defer func(consumer rabbitmq.Consumer) {
			err := consumer.Close()
			if err != nil {
				logger.Log.Error(err)
			}
		}(r.consumer)
	}
}

func NewRabbitMQConsumerImpl(rabbitMQConf conf.RabbitMQConf) *RabbitMQConsumerImpl {
	consumer, err := rabbitmq.NewConsumer(
		rabbitMQConf.GetURI(),
		rabbitmq.Config{},
		// can pass nothing for no logging
		rabbitmq.WithConsumerOptionsLogger(logger.Log),
	)
	if err != nil {
		logger.Log.Fatal(err)
	}
	return &RabbitMQConsumerImpl{consumer: consumer}
}
