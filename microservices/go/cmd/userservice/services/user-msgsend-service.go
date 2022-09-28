package services

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
	msg "github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq/exchanges"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq/rmqservices"
)

type UserMsgSendService interface {
	UserCreateSender() msg.Sender[userdtos.UserDto]
	UserUpdateSender() msg.Sender[userdtos.UserDto]
	UserDeleteSender() msg.Sender[userdtos.UserDto]
}

type UserMessageServiceImpl struct {
	userCreateSender msg.Sender[userdtos.UserDto]
	userUpdateSender msg.Sender[userdtos.UserDto]
	userDeleteSender msg.Sender[userdtos.UserDto]
}

func (u UserMessageServiceImpl) UserCreateSender() msg.Sender[userdtos.UserDto] {
	return u.userCreateSender
}

func (u UserMessageServiceImpl) UserUpdateSender() msg.Sender[userdtos.UserDto] {
	return u.userUpdateSender
}

func (u UserMessageServiceImpl) UserDeleteSender() msg.Sender[userdtos.UserDto] {
	return u.userDeleteSender
}

func NewUserMessageServiceImpl(connector rmqservices.RabbitConnector) UserMsgSendService {
	return &UserMessageServiceImpl{
		userCreateSender: rmq.NewSender(connector.GetPublisher(), exchanges.UserCreateExchange, []string{}),
		userUpdateSender: rmq.NewSender(connector.GetPublisher(), exchanges.UserUpdateExchange, []string{}),
		userDeleteSender: rmq.NewSender(connector.GetPublisher(), exchanges.UserDeleteExchange, []string{}),
	}
}
