package listeners

import (
	"context"
	"github.com/akrennmair/slice"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/lifecycle"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/kfka"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/kfka/topics"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
	"github.com/segmentio/kafka-go"
	"time"
)

type UserChangeListener interface {
	lifecycle.Closable
	ListenUserChange()
}

type UserChangeListenerImpl struct {
	userChangeEventService services.UserChangeEventService
	// Retry Topics
	user1Retry1Topic string
	user1Retry2Topic string
	user1Retry3Topic string
	user1Retry4Topic string
	// Dead letter Topic
	user1DeadLetterTopic string
	// Main receiver
	user1Receiver *kfka.KafkaReceiver[userdtos.UserChangeEventDto]
	// Retry receiver
	user1Retry1Receiver *kfka.KafkaReceiver[kfka.RetryDto[userdtos.UserChangeEventDto]]
	user1Retry2Receiver *kfka.KafkaReceiver[kfka.RetryDto[userdtos.UserChangeEventDto]]
	user1Retry3Receiver *kfka.KafkaReceiver[kfka.RetryDto[userdtos.UserChangeEventDto]]
	user1Retry4Receiver *kfka.KafkaReceiver[kfka.RetryDto[userdtos.UserChangeEventDto]]
	// Retry senders
	user1Retry1Sender     *kfka.KafkaSender[kfka.RetryDto[userdtos.UserChangeEventDto]]
	user1Retry2Sender     *kfka.KafkaSender[kfka.RetryDto[userdtos.UserChangeEventDto]]
	user1Retry3Sender     *kfka.KafkaSender[kfka.RetryDto[userdtos.UserChangeEventDto]]
	user1Retry4Sender     *kfka.KafkaSender[kfka.RetryDto[userdtos.UserChangeEventDto]]
	user1DeadLetterSender *kfka.KafkaSender[kfka.RetryDto[userdtos.UserChangeEventDto]]
}

func (k UserChangeListenerImpl) ListenUserChange() {
	k.user1Receiver.ListenSyncCommit(func(dto userdtos.UserChangeEventDto) error {
		ctx := context.Background()
		_, err := k.userChangeEventService.HandleUserChangeEventTxn(ctx, dto)
		if err != nil {
			logger.Log.WithError(err).Error("Error changing user")
			return k.user1Retry1Sender.Send(ctx, kfka.CreateRetryDto(dto, 0))
		}
		return nil
	})

	k.user1Retry1Receiver.ListenSyncCommit(func(retry kfka.RetryDto[userdtos.UserChangeEventDto]) error {
		ctx := context.Background()
		_, err := k.userChangeEventService.HandleUserChangeEventTxn(ctx, retry.Value)
		if err != nil {
			logger.Log.WithError(err).Error("Error changing user")
			if retry.Tries >= 3600 {
				return k.user1Retry2Sender.Send(ctx, retry)
			}
			retry.Tries += 1
			return k.user1Retry1Sender.Send(ctx, retry)
		}
		return nil
	})

	k.user1Retry2Receiver.ListenSyncCommit(func(retry kfka.RetryDto[userdtos.UserChangeEventDto]) error {
		ctx := context.Background()
		_, err := k.userChangeEventService.HandleUserChangeEventTxn(ctx, retry.Value)
		if err != nil {
			logger.Log.WithError(err).Error("Error changing user")
			if retry.Tries >= 1440 {
				return k.user1Retry3Sender.Send(ctx, retry)
			}
			retry.Tries += 1
			return k.user1Retry2Sender.Send(ctx, retry)
		}
		return nil
	})

	k.user1Retry3Receiver.ListenSyncCommit(func(retry kfka.RetryDto[userdtos.UserChangeEventDto]) error {
		ctx := context.Background()
		_, err := k.userChangeEventService.HandleUserChangeEventTxn(ctx, retry.Value)
		if err != nil {
			logger.Log.WithError(err).Error("Error changing user")
			if retry.Tries >= 100 {
				return k.user1Retry4Sender.Send(ctx, retry)
			}
			retry.Tries += 1
			return k.user1Retry3Sender.Send(ctx, retry)
		}
		return nil
	})

	k.user1Retry4Receiver.ListenSyncCommit(func(retry kfka.RetryDto[userdtos.UserChangeEventDto]) error {
		ctx := context.Background()
		_, err := k.userChangeEventService.HandleUserChangeEventTxn(ctx, retry.Value)
		if err != nil {
			logger.Log.WithError(err).Error("Error changing user")
			if retry.Tries >= 20 {
				return k.user1DeadLetterSender.Send(ctx, retry)
			}
			retry.Tries += 1
			return k.user1Retry4Sender.Send(ctx, retry)
		}
		return nil
	})
	logger.Log.Info("Listening for user changes")
}

func (k UserChangeListenerImpl) Close() error {
	logger.Log.Info("Closing user listener")
	closableList := []lifecycle.Closable{
		k.user1Receiver,
		k.user1Retry1Receiver,
		k.user1Retry2Receiver,
		k.user1Retry3Receiver,
		k.user1Retry4Receiver,
		k.user1Retry1Sender,
		k.user1Retry2Sender,
		k.user1Retry3Sender,
		k.user1Retry4Sender,
		k.user1DeadLetterSender,
	}
	errs := slice.MapConcurrent(closableList, lifecycle.Closable.Close)
	err := utils.ConcatErrors(errs...)
	if err != nil {
		logger.Log.WithError(err).Error("Errors closing listener")
	}
	logger.Log.Info("Finish closing user listener")
	return err
}

func NewUserListenerImpl(
	userChangeEventService services.UserChangeEventService,
	kafkaConf conf.KafkaConf,
) *UserChangeListenerImpl {
	if !environment.ActivateKafkaListener() {
		// Listener is deactivated, ran via the lifecycle package,
		// and is a root-child dependency so a nil is returned
		return nil
	}

	serviceAppend := "-key-service"

	user1Retry1Topic := topics.AppendRetry(topics.User1Topic+serviceAppend, 1)
	user1Retry2Topic := topics.AppendRetry(topics.User1Topic+serviceAppend, 2)
	user1Retry3Topic := topics.AppendRetry(topics.User1Topic+serviceAppend, 3)
	user1Retry4Topic := topics.AppendRetry(topics.User1Topic+serviceAppend, 4)
	user1DeadLetterTopic := topics.AppendDeadLetter(topics.User1Topic + serviceAppend)

	user1Receiver := kfka.NewKafkaReceiver[userdtos.UserChangeEventDto](
		kafka.NewReader(kafka.ReaderConfig{
			Brokers:  kafkaConf.GetBootstrapServers(),
			GroupID:  topics.User1Topic + serviceAppend,
			Topic:    topics.User1Topic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
	)
	user1Retry1Receiver := kfka.NewKafkaReceiver[kfka.RetryDto[userdtos.UserChangeEventDto]](
		kafka.NewReader(kafka.ReaderConfig{
			Brokers:        kafkaConf.GetBootstrapServers(),
			GroupID:        user1Retry1Topic,
			Topic:          user1Retry1Topic,
			MinBytes:       10e3, // 10KB
			MaxBytes:       10e6, // 10MB,
			ReadBackoffMin: time.Second,
			ReadBackoffMax: 10 * time.Second,
		}),
	)
	user1Retry2Receiver := kfka.NewKafkaReceiver[kfka.RetryDto[userdtos.UserChangeEventDto]](
		kafka.NewReader(kafka.ReaderConfig{
			Brokers:        kafkaConf.GetBootstrapServers(),
			GroupID:        user1Retry2Topic,
			Topic:          user1Retry2Topic,
			MinBytes:       10e3, // 10KB
			MaxBytes:       10e6, // 10MB
			ReadBackoffMin: time.Minute,
			ReadBackoffMax: 2 * time.Minute,
		}),
	)
	user1Retry3Receiver := kfka.NewKafkaReceiver[kfka.RetryDto[userdtos.UserChangeEventDto]](
		kafka.NewReader(kafka.ReaderConfig{
			Brokers:        kafkaConf.GetBootstrapServers(),
			GroupID:        user1Retry3Topic,
			Topic:          user1Retry3Topic,
			MinBytes:       10e3, // 10KB
			MaxBytes:       10e6, // 10MB
			ReadBackoffMin: time.Hour,
			ReadBackoffMax: time.Hour + 15*time.Minute,
		}),
	)
	user1Retry4Receiver := kfka.NewKafkaReceiver[kfka.RetryDto[userdtos.UserChangeEventDto]](
		kafka.NewReader(kafka.ReaderConfig{
			Brokers:        kafkaConf.GetBootstrapServers(),
			GroupID:        user1Retry4Topic,
			Topic:          user1Retry4Topic,
			MinBytes:       10e3, // 10KB
			MaxBytes:       10e6, // 10MB
			ReadBackoffMin: 5 * time.Hour,
			ReadBackoffMax: 10 * time.Hour,
		}),
	)
	user1Retry1Sender := kfka.NewKafkaSender(
		&kafka.Writer{
			Addr:     kafka.TCP(kafkaConf.GetBootstrapServers()...),
			Topic:    user1Retry1Topic,
			Balancer: &kafka.Murmur2Balancer{},
		},
		func(b kfka.RetryDto[userdtos.UserChangeEventDto]) ([]byte, error) {
			return b.Value.MessageKey()
		},
	)
	user1Retry2Sender := kfka.NewKafkaSender(
		&kafka.Writer{
			Addr:     kafka.TCP(kafkaConf.GetBootstrapServers()...),
			Topic:    user1Retry2Topic,
			Balancer: &kafka.Murmur2Balancer{},
		},
		func(b kfka.RetryDto[userdtos.UserChangeEventDto]) ([]byte, error) {
			return b.Value.MessageKey()
		},
	)
	user1Retry3Sender := kfka.NewKafkaSender(
		&kafka.Writer{
			Addr:     kafka.TCP(kafkaConf.GetBootstrapServers()...),
			Topic:    user1Retry3Topic,
			Balancer: &kafka.Murmur2Balancer{},
		},
		func(b kfka.RetryDto[userdtos.UserChangeEventDto]) ([]byte, error) {
			return b.Value.MessageKey()
		},
	)
	user1Retry4Sender := kfka.NewKafkaSender(
		&kafka.Writer{
			Addr:     kafka.TCP(kafkaConf.GetBootstrapServers()...),
			Topic:    user1Retry4Topic,
			Balancer: &kafka.Murmur2Balancer{},
		},
		func(b kfka.RetryDto[userdtos.UserChangeEventDto]) ([]byte, error) {
			return b.Value.MessageKey()
		},
	)
	user1DeadLetterSender := kfka.NewKafkaSender(
		&kafka.Writer{
			Addr:     kafka.TCP(kafkaConf.GetBootstrapServers()...),
			Topic:    user1DeadLetterTopic,
			Balancer: &kafka.Murmur2Balancer{},
		},
		func(b kfka.RetryDto[userdtos.UserChangeEventDto]) ([]byte, error) {
			return b.Value.MessageKey()
		},
	)
	r := &UserChangeListenerImpl{
		userChangeEventService: userChangeEventService,
		// Topics
		user1Retry1Topic:     user1Retry1Topic,
		user1Retry2Topic:     user1Retry2Topic,
		user1Retry3Topic:     user1Retry3Topic,
		user1Retry4Topic:     user1Retry4Topic,
		user1DeadLetterTopic: user1DeadLetterTopic,
		// Receiver
		user1Receiver:       user1Receiver,
		user1Retry1Receiver: user1Retry1Receiver,
		user1Retry2Receiver: user1Retry2Receiver,
		user1Retry3Receiver: user1Retry3Receiver,
		user1Retry4Receiver: user1Retry4Receiver,
		// Sender
		user1Retry1Sender:     user1Retry1Sender,
		user1Retry2Sender:     user1Retry2Sender,
		user1Retry3Sender:     user1Retry3Sender,
		user1Retry4Sender:     user1Retry4Sender,
		user1DeadLetterSender: user1DeadLetterSender,
	}
	lifecycle.RegisterClosable(r)
	return r
}
