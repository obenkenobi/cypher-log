package services

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq/exchanges"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
)

type UserMessageService interface {
	SendUserCreate(userDto userdtos.UserDto) single.Single[userdtos.UserDto]
	SendUserUpdate(userDto userdtos.UserDto) single.Single[userdtos.UserDto]
	SendUserDelete(userDto userdtos.UserDto) single.Single[userdtos.UserDto]
}

type UserMessageServiceImpl struct {
	rabbitMQConf     conf.RabbitMQConf
	userCreateSender messaging.Sender[userdtos.UserDto]
	userUpdateSender messaging.Sender[userdtos.UserDto]
	userDeleteSender messaging.Sender[userdtos.UserDto]
}

func (u UserMessageServiceImpl) SendUserCreate(userDto userdtos.UserDto) single.Single[userdtos.UserDto] {
	return u.userCreateSender.Send(userDto)
}

func (u UserMessageServiceImpl) SendUserUpdate(userDto userdtos.UserDto) single.Single[userdtos.UserDto] {
	return u.userUpdateSender.Send(userDto)
}

func (u UserMessageServiceImpl) SendUserDelete(userDto userdtos.UserDto) single.Single[userdtos.UserDto] {
	return u.userDeleteSender.Send(userDto)
}

func NewUserMessageServiceImpl(rabbitMQConf conf.RabbitMQConf) UserMessageService {
	return &UserMessageServiceImpl{
		rabbitMQConf: rabbitMQConf,
		userCreateSender: rmq.NewRabbitMQSender(rmq.SenderOptions[userdtos.UserDto]{
			ExchangeOpts: exchanges.UserCreateExchangeOpts,
			RoutingKey:   "",
			Mandatory:    false,
			Immediate:    false,
		}),
		userUpdateSender: rmq.NewRabbitMQSender(rmq.SenderOptions[userdtos.UserDto]{
			ExchangeOpts: exchanges.UserUpdateExchangeOpts,
			RoutingKey:   "",
			Mandatory:    false,
			Immediate:    false,
		}),
		userDeleteSender: rmq.NewRabbitMQSender(rmq.SenderOptions[userdtos.UserDto]{
			ExchangeOpts: exchanges.UserDeleteExchangeOpts,
			RoutingKey:   "",
			Mandatory:    false,
			Immediate:    false,
		}),
	}
}
