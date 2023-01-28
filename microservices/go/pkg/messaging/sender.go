package messaging

type Sender[T any] interface {
	Send(body T) error
}
