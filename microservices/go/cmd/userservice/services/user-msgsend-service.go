package services

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
	msg "github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq/exchanges"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/rmqservices"
)

type UserMsgSendService interface {
	UserSaveSender() msg.Sender[userdtos.DistributedUserDto]
	UserDeleteSender() msg.Sender[userdtos.DistributedUserDto]
}

type UserMessageServiceImpl struct {
	userSaveSender   msg.Sender[userdtos.DistributedUserDto]
	userDeleteSender msg.Sender[userdtos.DistributedUserDto]
}

func (u UserMessageServiceImpl) UserSaveSender() msg.Sender[userdtos.DistributedUserDto] {
	return u.userSaveSender
}

func (u UserMessageServiceImpl) UserDeleteSender() msg.Sender[userdtos.DistributedUserDto] {
	return u.userDeleteSender
}

func NewUserMessageServiceImpl(connector rmqservices.RabbitConnector) UserMsgSendService {
	return &UserMessageServiceImpl{
		userSaveSender:   rmq.NewSender(connector.GetPublisher(), exchanges.UserSaveExchange, rmq.RoutingKeysDefault),
		userDeleteSender: rmq.NewSender(connector.GetPublisher(), exchanges.UserDeleteExchange, rmq.RoutingKeysDefault),
	}
}
