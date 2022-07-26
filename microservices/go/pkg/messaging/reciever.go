package messaging

// Receiver response

type ReceiverAction int

const (
	Commit ReceiverAction = iota
	Discard
	Resend
)

// Receiver message

type Delivery[T any] interface {
	Body() T
	Commit() ReceiverAction
	Discard() ReceiverAction
	Resend() ReceiverAction
}

type receiverMessage[T any] struct {
	body T
}

func (r receiverMessage[T]) Body() T {
	return r.body
}

func (r receiverMessage[T]) Commit() ReceiverAction {
	return Commit
}

func (r receiverMessage[T]) Discard() ReceiverAction {
	return Discard
}

func (r receiverMessage[T]) Resend() ReceiverAction {
	return Resend
}

func NewDelivery[T any](body T) Delivery[T] {
	return receiverMessage[T]{body: body}
}

type Receiver[T any] interface {
	Listen(listener func(delivery Delivery[T]) ReceiverAction)
}
