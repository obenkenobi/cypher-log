package exchanges

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/messaging/rmq"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/userdtos"
)

var UserChangeExchange = rmq.Exchange[userdtos.UserChangeEventDto]{
	Name:        "user_change",
	Kind:        rmq.ExchangeTypeFanout,
	Durable:     true,
	AutoDeleted: false,
	Internal:    false,
	NoWait:      false,
}
