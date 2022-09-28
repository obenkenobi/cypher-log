package exchanges

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq"
)

var UserSaveExchange = rmq.Exchange[userdtos.UserDto]{
	Name:        "user_save",
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
