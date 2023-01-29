package listeners

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/lifecycle"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/kfka"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/userdtos"
	"github.com/segmentio/kafka-go"
)

type KafkaListener interface {
	lifecycle.TaskRunner
	lifecycle.Closable
}

type KafkaListenerImpl struct {
	userChangeEventService services.UserChangeEventService
	userChangeReceiver     *kfka.KafkaReceiver[userdtos.UserChangeEventDto]
}

func (k KafkaListenerImpl) ListenUserChange() {
	k.userChangeReceiver.Listen(func(dto userdtos.UserChangeEventDto) error {
		ctx := context.Background()
		_, err := k.userChangeEventService.HandleUserChangeEventTxn(ctx, dto)
		return err
	})
	logger.Log.Info("Listening for user changes")
}

func (k KafkaListenerImpl) Run() {
	k.ListenUserChange()
	forever := make(chan any)
	<-forever
}

func (k KafkaListenerImpl) Close() error {
	err := k.userChangeReceiver.Close()
	if err != nil {
		logger.Log.WithError(err).Error()
	}
	return err
}

func NewKafkaListenerImpl(
	userChangeEventService services.UserChangeEventService,
	kafkaConf conf.KafkaConf,
) *KafkaListenerImpl {
	if !environment.ActivateKafkaListener() {
		// Listener is deactivated, ran via the lifecycle package,
		// and is a root-child dependency so a nil is returned
		return nil
	}

	userChangeReceiver := kfka.NewKafkaReceiver[userdtos.UserChangeEventDto](
		kafka.NewReader(kafka.ReaderConfig{
			Brokers:  kafkaConf.GetServers(),
			GroupID:  "user-0-note-service-0",
			Topic:    "user-0",
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
	)
	r := &KafkaListenerImpl{
		userChangeReceiver:     userChangeReceiver,
		userChangeEventService: userChangeEventService,
	}
	lifecycle.RegisterClosable(r)
	lifecycle.RegisterTaskRunner(r)
	return r
}
