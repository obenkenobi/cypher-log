package exchanges

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq"
)

var UserCreateExchangeOpts = rmq.ExchangeOptions[userdtos.UserDto]{
	Name:        "user_create",
	Kind:        rmq.ExchangeTypeFanout,
	Durable:     true,
	AutoDeleted: false,
	Internal:    false,
	NoWait:      false,
	Args:        nil,
}

var UserUpdateExchangeOpts = rmq.ExchangeOptions[userdtos.UserDto]{
	Name:        "user_update",
	Kind:        rmq.ExchangeTypeFanout,
	Durable:     true,
	AutoDeleted: false,
	Internal:    false,
	NoWait:      false,
	Args:        nil,
}

var UserDeleteExchangeOpts = rmq.ExchangeOptions[userdtos.UserDto]{
	Name:        "user_delete",
	Kind:        rmq.ExchangeTypeFanout,
	Durable:     true,
	AutoDeleted: false,
	Internal:    false,
	NoWait:      false,
	Args:        nil,
}
