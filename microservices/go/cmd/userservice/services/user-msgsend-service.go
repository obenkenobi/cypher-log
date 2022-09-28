package services

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
	msg "github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq/exchanges"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq/rmqservices"
)

type UserMsgSendService interface {
	UserSaveSender() msg.Sender[userdtos.UserDto]
	UserDeleteSender() msg.Sender[userdtos.UserDto]
}

type UserMessageServiceImpl struct {
	userSaveSender   msg.Sender[userdtos.UserDto]
	userDeleteSender msg.Sender[userdtos.UserDto]
}

func (u UserMessageServiceImpl) UserSaveSender() msg.Sender[userdtos.UserDto] {
	return u.userSaveSender
}

func (u UserMessageServiceImpl) UserDeleteSender() msg.Sender[userdtos.UserDto] {
	return u.userDeleteSender
}

func NewUserMessageServiceImpl(connector rmqservices.RabbitConnector) UserMsgSendService {
	return &UserMessageServiceImpl{
		userSaveSender:   rmq.NewSender(connector.GetPublisher(), exchanges.UserSaveExchange, rmq.RoutingKeysDefault),
		userDeleteSender: rmq.NewSender(connector.GetPublisher(), exchanges.UserDeleteExchange, rmq.RoutingKeysDefault),
	}
}
