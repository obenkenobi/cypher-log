package kfka

import (
	"context"
	"encoding/json"
	"github.com/barweiss/go-tuple"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/lifecycle"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
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

// ListenSyncCommit listens for messages and commits after processMsg is
// finished. If processMsg returns an error, the listener ends and exits the
// application. If the intention of the processing is to continue even if there
// is an error, return a nil.
func (r *KafkaReceiver[T]) ListenSyncCommit(processMsg func(T) error) bool {
	return r.Listen(false, processMsg)
}

// ListenAutoCommit listens for messages and commits the moment the message is
// read. If processMsg returns an error, the listener ends and exits the
// application. If the intention of the processing is to continue even if there
// is an error, return a nil.
func (r *KafkaReceiver[T]) ListenAutoCommit(processMsg func(T) error) bool {
	return r.Listen(true, processMsg)
}

// Listen listens for messages. Setting autoCommit to true commits the message
// the moment it is read. Setting autoCommit to false commits the message after
// the message is processed. If processMsg returns an error, the listener ends
// and exits the application. If the intention of the processing is to continue
// even if there is an error, return a nil.
func (r *KafkaReceiver[T]) Listen(autoCommit bool, processMsg func(T) error) bool {
	if r.shouldNotListen() {
		return false
	}
	go r.runListen(autoCommit, processMsg)
	return true
}

func (r *KafkaReceiver[T]) Close() error {
	r.closeMutex.Lock()
	defer r.closeMutex.Unlock()
	if r.isClosed {
		return nil
	}
	go func() {
		r.beginCloseCh <- true
	}()
	err := <-r.closedCh
	r.isClosed = true
	return err
}

func (r *KafkaReceiver[T]) runListen(autoCommit bool, processMsg func(T) error) {
	msgChan := make(chan tuple.T2[context.Context, kafka.Message])
	errChan := make(chan tuple.T2[context.Context, error])

	go func() {
		for {
			ctx := context.Background()
			var m kafka.Message
			var err error
			if autoCommit {
				m, err = r.reader.ReadMessage(ctx)
			} else {
				m, err = r.reader.FetchMessage(ctx)
			}
			if err != nil {
				errChan <- tuple.New2(ctx, err)
			}
			msgChan <- tuple.New2(ctx, m)
		}
	}()

	var shutDownErr error = nil
	willContinue := true
	for willContinue {
		select {
		case msgTuple := <-msgChan:
			ctx, msg := msgTuple.V1, msgTuple.V2
			if body, err := r.readBody(msg); err != nil {
				r.logger().WithContext(ctx).WithError(err).Error("Cannot parse message value, skipping message")
			} else {
				if err := processMsg(body); err != nil {
					r.logger().WithContext(ctx).WithError(err).Error("Failed to process message")
					willContinue = false
					shutDownErr = err
					continue
				}
			}
			if !autoCommit {
				if err := r.reader.CommitMessages(ctx, msg); err != nil {
					r.logger().WithContext(ctx).WithError(err).Error("Error committing a message")
					willContinue = false
					shutDownErr = err
					continue
				}
			}
		case errTuple := <-errChan:
			ctx, err := errTuple.V1, errTuple.V2
			r.logger().WithContext(ctx).WithError(err).Error("Error reading a message")
			willContinue = false
			shutDownErr = err
		case <-r.beginCloseCh:
			willContinue = false
		}
	}
	go func() {
		r.closedCh <- r.reader.Close()
	}()

	if shutDownErr != nil {
		r.logger().
			WithError(shutDownErr).
			Errorf("Ending receiver for topic %v due to an error", r.reader.Config().Topic)
		lifecycle.ExitApp()
	} else {
		r.logger().Infof("Ending receiver for topic %v", r.reader.Config().Topic)
	}
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

func (r *KafkaReceiver[T]) logger() *log.Entry {
	return logger.Log.WithField("receiverTopic", r.reader.Config().Topic)
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
