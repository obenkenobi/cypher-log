package messaging

type Receiver[T any] interface {
	Listen(listener func(body T) error, resendIfErr bool)
}
