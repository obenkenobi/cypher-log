package lifecycle

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
	"sync"
)

// Closable is an interface for an object that has closable resources
type Closable interface {
	Close() error
}

var closableList []Closable
var closableMutex sync.Mutex

// RegisterClosable registers a Closable in your app lifecycle to be closed gracefully
func RegisterClosable(closable Closable) {
	closableMutex.Lock()
	closableMutex.Unlock()
	closableList = append(closableList, closable)
}

func closeResources() bool {
	isGraceful := true
	for _, closable := range closableList {
		logger.Log.WithField("Type", utils.GetType(closable)).Info("Closing a resource")
		if err := closable.Close(); err != nil {
			logger.Log.WithError(err).Error("Error while closing")
			isGraceful = false
		}
	}
	return isGraceful
}
