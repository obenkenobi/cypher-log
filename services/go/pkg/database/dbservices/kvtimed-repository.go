package dbservices

import (
	"context"
	"github.com/obenkenobi/cypher-log/services/go/pkg/reactive/single"
)

type KeyValueTimedRepository[Key any, Value any] interface {
	Get(ctx context.Context, key Key) single.Single[Value]
	Set(ctx context.Context, key Key, value Value) single.Single[Value]
}

// Todo: create a redis implementation of KeyValueTimedRepository
