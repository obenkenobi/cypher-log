package listeners

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/lifecycle"
)

type KafkaListener interface {
	lifecycle.TaskRunner
}

type KafkaListenerImpl struct {
	userListener UserListener
}

func (k KafkaListenerImpl) Run() {
	k.userListener.ListenUserChange()
	forever := make(chan any)
	<-forever
}

func NewKafkaListenerImpl(
	userListener UserListener,
) *KafkaListenerImpl {
	if !environment.ActivateKafkaListener() {
		// Listener is deactivated, ran via the lifecycle package,
		// and is a root-child dependency so a nil is returned
		return nil
	}
	r := &KafkaListenerImpl{userListener: userListener}
	lifecycle.RegisterTaskRunner(r)
	return r
}
