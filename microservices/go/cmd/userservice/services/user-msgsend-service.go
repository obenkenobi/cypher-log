package services

import (
	msg "github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq/exchanges"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/rmqservices"
)

type UserMsgSendService interface {
	SendUserSave(dto userdtos.UserChangeEventDto) error
}

type UserMessageServiceImpl struct {
	userSaveSender msg.Sender[userdtos.UserChangeEventDto]
}

func (u UserMessageServiceImpl) SendUserSave(dto userdtos.UserChangeEventDto) error {
	_, err := u.userSaveSender.Send(dto)
	return err
}

func NewUserMessageServiceImpl(publisher rmqservices.RabbitMQPublisher) *UserMessageServiceImpl {
	return &UserMessageServiceImpl{
		userSaveSender: rmq.NewSender(publisher.GetPublisher(), exchanges.UserChangeExchange, rmq.RoutingKeysDefault),
	}
}
