package rmqservices

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/wagslane/go-rabbitmq"
)

type RabbitConnector interface {
	GetConsumer() rabbitmq.Consumer
	GetPublisher() *rabbitmq.Publisher
	Close()
}

type rabbitConnectorImpl struct {
	publisher *rabbitmq.Publisher
	consumer  rabbitmq.Consumer
}

func (r rabbitConnectorImpl) GetConsumer() rabbitmq.Consumer { return r.consumer }

func (r rabbitConnectorImpl) GetPublisher() *rabbitmq.Publisher { return r.publisher }

func (r *rabbitConnectorImpl) Close() {
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
		defer func(consumer rabbitmq.Consumer) {
			err := consumer.Close()
			if err != nil {
				logger.Log.Error(err)
			}
		}(r.consumer)

	}
}

func NewRabbitConnector(rabbitMQConf conf.RabbitMQConf) RabbitConnector {
	publisher, err := rabbitmq.NewPublisher(
		rabbitMQConf.GetURI(),
		rabbitmq.Config{},
		// can pass nothing for no logging
		rabbitmq.WithPublisherOptionsLogger(logger.Log),
	)
	if err != nil {
		logger.Log.Fatal(err)
	}
	consumer, err := rabbitmq.NewConsumer(
		rabbitMQConf.GetURI(),
		rabbitmq.Config{},
		// can pass nothing for no logging
		rabbitmq.WithConsumerOptionsLogger(logger.Log),
	)
	if err != nil {
		logger.Log.Fatal(err)
	}
	return &rabbitConnectorImpl{publisher: publisher, consumer: consumer}
}
