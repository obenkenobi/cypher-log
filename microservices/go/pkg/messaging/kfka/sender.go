package kfka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"sync"
)

// KafkaSender is a thread safe sender class to send messages over kafka
type KafkaSender[T any] struct {
	writer    *kafka.Writer
	rwLock    sync.RWMutex
	keyReader func(body T) ([]byte, error)
}

func (k KafkaSender[T]) Send(ctx context.Context, body T) error {
	k.rwLock.RLock()
	k.rwLock.RUnlock()
	msgBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	key, err := k.keyReader(body)
	if err != nil {
		return err
	}
	return k.writer.WriteMessages(ctx, kafka.Message{Key: key, Value: msgBytes})
}

func (k KafkaSender[T]) Close() error {
	k.rwLock.Lock()
	k.rwLock.Unlock()
	if k.writer != nil {
		return k.writer.Close()
	}
	return nil
}

func NewKafkaSender[T any](writer *kafka.Writer, keyReader func(body T) ([]byte, error)) *KafkaSender[T] {
	return &KafkaSender[T]{writer: writer, keyReader: keyReader}
}
