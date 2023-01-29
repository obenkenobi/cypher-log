package services

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/lifecycle"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/kfka"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/userdtos"
	"github.com/segmentio/kafka-go"
)

type UserMsgSendService interface {
	SendUserSave(ctx context.Context, dto userdtos.UserChangeEventDto) error
	lifecycle.Closable
}

type UserMessageServiceImpl struct {
	userSaveSender kfka.KafkaSender[userdtos.UserChangeEventDto]
}

func (u *UserMessageServiceImpl) SendUserSave(ctx context.Context, dto userdtos.UserChangeEventDto) error {
	return u.userSaveSender.Send(ctx, dto)
}

func (u *UserMessageServiceImpl) Close() error {
	return u.userSaveSender.Close()
}

func NewUserMessageServiceImpl(kafkaConf conf.KafkaConf) *UserMessageServiceImpl {
	userSaveSender := kfka.NewKafkaSender(
		&kafka.Writer{
			Addr:     kafka.TCP(kafkaConf.GetBootstrapServers()...),
			Topic:    "user-0",
			Balancer: &kafka.Murmur2Balancer{},
		},
		func(dto userdtos.UserChangeEventDto) ([]byte, error) {
			return []byte(dto.Id), nil
		})
	u := &UserMessageServiceImpl{userSaveSender: userSaveSender}
	lifecycle.RegisterClosable(u)
	return u
}
