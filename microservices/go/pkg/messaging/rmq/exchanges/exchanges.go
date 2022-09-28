package exchanges

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq"
)

var UserCreateExchange = rmq.Exchange[userdtos.UserDto]{
	Name:        "user_create",
	Kind:        rmq.ExchangeTypeFanout,
	Durable:     true,
	AutoDeleted: false,
	Internal:    false,
	NoWait:      false,
}

var UserUpdateExchange = rmq.Exchange[userdtos.UserDto]{
	Name:        "user_update",
	Kind:        rmq.ExchangeTypeFanout,
	Durable:     true,
	AutoDeleted: false,
	Internal:    false,
	NoWait:      false,
}

var UserDeleteExchange = rmq.Exchange[userdtos.UserDto]{
	Name:        "user_delete",
	Kind:        rmq.ExchangeTypeFanout,
	Durable:     true,
	AutoDeleted: false,
	Internal:    false,
	NoWait:      false,
}
