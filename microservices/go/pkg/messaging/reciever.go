package messaging

import "github.com/joamaki/goreactive/stream"

type Receiver[T any] interface {
	ReceiveMessages() stream.Observable[T]
	IsAutoAck() bool
	Close()
}

type Delivered[T any] interface {
	GetBody() T
	Ack(multiple bool) error
	Reject(sendBack bool) error
	IsAutoAck() bool
}
