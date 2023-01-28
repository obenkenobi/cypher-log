package rmq

import (
	"encoding/json"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging"
	"github.com/wagslane/go-rabbitmq"
)

type Sender[T any] struct {
	publisher   *rabbitmq.Publisher
	routingKeys []string
	publishOpts []func(options *rabbitmq.PublishOptions)
}

func NewSender[T any](
	publisher *rabbitmq.Publisher,
	exchange Exchange[T],
	routingKeys []string,
	publishOpts ...func(options *rabbitmq.PublishOptions),
) messaging.Sender[T] {
	publishOpts = append(publishOpts, exchange.GetPublishOptions()...)
	return &Sender[T]{
		publisher:   publisher,
		routingKeys: routingKeys,
		publishOpts: publishOpts,
	}
}

func (r *Sender[T]) Send(body T) error {
	var msgBytes []byte
	var contentType string
	{
		var bodyMatcher interface{} = body
		switch v := bodyMatcher.(type) {
		case string:
			msgBytes = []byte(v)
			contentType = ContentTypePlainText
		default:
			var err error
			if msgBytes, err = json.Marshal(body); err != nil {
				return err
			}
			contentType = ContentTypeJson
		}
	}
	publishOpts := append(r.publishOpts,
		rabbitmq.WithPublishOptionsContentType(contentType),
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsMandatory)
	r.publisher.NotifyReturn()
	err := r.publisher.Publish(msgBytes, r.routingKeys, publishOpts...)
	return err
}
