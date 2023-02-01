package lifecycle

import (
	"github.com/akrennmair/slice"
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
	errChannels := slice.Map(closableList, func(c Closable) chan error {
		errC := make(chan error)
		go func(closable Closable, errCh chan error) {
			errCh <- closable.Close()
		}(c, errC)
		return errC
	})
	isGraceful := true
	for i, errCh := range errChannels {
		if err := <-errCh; err != nil {
			logger.Log.WithError(err).Error("Error while closing")
			isGraceful = false
		} else {
			closableType := utils.GetType(closableList[i])
			logger.Log.WithField("ClosableType", closableType).Infof("Closed %v", closableType)
		}
	}
	return isGraceful
}
