package services

import (
	msg "github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq/exchanges"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/rmqservices"
)

type UserMsgSendService interface {
	UserSaveSender() msg.Sender[userdtos.DistUserSaveDto]
	UserDeleteSender() msg.Sender[userdtos.DistUserDeleteDto]
}

type UserMessageServiceImpl struct {
	userSaveSender   msg.Sender[userdtos.DistUserSaveDto]
	userDeleteSender msg.Sender[userdtos.DistUserDeleteDto]
}

func (u UserMessageServiceImpl) UserSaveSender() msg.Sender[userdtos.DistUserSaveDto] {
	return u.userSaveSender
}

func (u UserMessageServiceImpl) UserDeleteSender() msg.Sender[userdtos.DistUserDeleteDto] {
	return u.userDeleteSender
}

func NewUserMessageServiceImpl(connector rmqservices.RabbitConnector) UserMsgSendService {
	return &UserMessageServiceImpl{
		userSaveSender:   rmq.NewSender(connector.GetPublisher(), exchanges.UserSaveExchange, rmq.RoutingKeysDefault),
		userDeleteSender: rmq.NewSender(connector.GetPublisher(), exchanges.UserDeleteExchange, rmq.RoutingKeysDefault),
	}
}
