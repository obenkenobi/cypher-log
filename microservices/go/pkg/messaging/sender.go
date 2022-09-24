package messaging

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
)

type Sender[T any] interface {
	Send(T) single.Single[T]
}
