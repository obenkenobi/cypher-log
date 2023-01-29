package kfka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/barweiss/go-tuple"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/segmentio/kafka-go"
	"io"
	"sync"
)

// KafkaReceiver is a thread safe solution to a kafka consumer
type KafkaReceiver[T any] struct {
	reader       *kafka.Reader
	isListening  bool
	listenMutex  sync.Mutex // using a mutex instead of a rw lock because the receiver should only listen once
	isClosed     bool
	closeMutex   sync.Mutex // using a mutex instead of a rw lock because the receiver should only close once
	closedCh     chan error
	beginCloseCh chan bool
}

func (r *KafkaReceiver[T]) Listen(handler func(T) error) bool {
	if r.shouldNotListen() {
		return false
	}
	go r.runListen(handler)
	return true
}

func (r *KafkaReceiver[T]) Close() error {
	r.closeMutex.Lock()
	defer r.closeMutex.Unlock()
	if r.isClosed {
		return nil
	}
	r.beginCloseCh <- true
	err := <-r.closedCh
	r.isClosed = true
	return err
}

func (r *KafkaReceiver[T]) runListen(handler func(T) error) {
	msgChan := make(chan tuple.T2[context.Context, kafka.Message])
	errChan := make(chan tuple.T2[context.Context, error])

	go func() {
		for {
			ctx := context.Background()
			m, err := r.reader.FetchMessage(ctx)
			if err != nil {
				errChan <- tuple.New2(ctx, err)
			}
			msgChan <- tuple.New2(ctx, m)
		}
	}()

	willContinue := true
	for willContinue {
		select {
		case msgTuple := <-msgChan:
			_, msg := msgTuple.V1, msgTuple.V2
			body, err := r.readBody(msg)
			if err != nil {
				logger.Log.WithError(err).Error()
				continue
			}
			err = handler(body)
			if err != nil {
				logger.Log.WithError(err).Error()
				continue
			}
		case errTuple := <-errChan:
			_, err := errTuple.V1, errTuple.V2
			logger.Log.WithError(err).Error()
			if errors.Is(err, io.EOF) {
				willContinue = false
			}
		case <-r.beginCloseCh:
			willContinue = false
		}
	}
	r.closedCh <- r.reader.Close()
}

func (r *KafkaReceiver[T]) shouldNotListen() bool {
	r.listenMutex.Lock()
	defer r.listenMutex.Unlock()
	if r.isListening {
		return true
	}
	r.isListening = true
	return false
}

func (r *KafkaReceiver[T]) readBody(msg kafka.Message) (T, error) {
	var body T
	err := json.Unmarshal(msg.Value, &body)
	return body, err
}

func NewKafkaReceiver[T any](reader *kafka.Reader) *KafkaReceiver[T] {
	return &KafkaReceiver[T]{
		reader:       reader,
		isListening:  false,
		isClosed:     false,
		beginCloseCh: make(chan bool),
		closedCh:     make(chan error),
	}
}
