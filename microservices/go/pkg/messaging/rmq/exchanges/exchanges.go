package exchanges

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/userdtos"
)

var UserSaveExchange = rmq.Exchange[userdtos.DistUserSaveDto]{
	Name:        "user_save",
	Kind:        rmq.ExchangeTypeFanout,
	Durable:     true,
	AutoDeleted: false,
	Internal:    false,
	NoWait:      false,
}

var UserDeleteExchange = rmq.Exchange[userdtos.DistUserDeleteDto]{
	Name:        "user_delete",
	Kind:        rmq.ExchangeTypeFanout,
	Durable:     true,
	AutoDeleted: false,
	Internal:    false,
	NoWait:      false,
}
