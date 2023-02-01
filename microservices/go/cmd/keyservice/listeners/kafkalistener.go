package listeners

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/lifecycle"
)

type KafkaListener interface {
	lifecycle.TaskRunner
}

type KafkaListenerImpl struct {
	userChange1Listener UserChange1Listener
}

func (k KafkaListenerImpl) Run() {
	k.userChange1Listener.ListenUserChange()
	forever := make(chan any)
	<-forever
}

func NewKafkaListenerImpl(
	userListener UserChange1Listener,
) *KafkaListenerImpl {
	if !environment.ActivateKafkaListener() {
		// Listener is deactivated, ran via the lifecycle package,
		// and is a root-child dependency so a nil is returned
		return nil
	}
	r := &KafkaListenerImpl{userChange1Listener: userListener}
	lifecycle.RegisterTaskRunner(r)
	return r
}
