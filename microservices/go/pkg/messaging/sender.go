package messaging

import "context"

type Sender[T any] interface {
	Send(ctx context.Context, body T) error
	Close() error
}
