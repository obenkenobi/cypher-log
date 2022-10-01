package services

import (
	msg "github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq/exchanges"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/rmqservices"
)

type UserMsgSendService interface {
	UserSaveSender() msg.Sender[userdtos.UserChangeEventDto]
}

type UserMessageServiceImpl struct {
	userSaveSender msg.Sender[userdtos.UserChangeEventDto]
}

func (u UserMessageServiceImpl) UserSaveSender() msg.Sender[userdtos.UserChangeEventDto] {
	return u.userSaveSender
}

func NewUserMessageServiceImpl(connector rmqservices.RabbitConnector) UserMsgSendService {
	return &UserMessageServiceImpl{
		userSaveSender: rmq.NewSender(connector.GetPublisher(), exchanges.UserChangeExchange, rmq.RoutingKeysDefault),
	}
}
