package lifecycle

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"os"
)

// RunApp starts your application concurrently by running each task in
// taskRunners concurrently and closes registered resources when the app ends
func RunApp() {
	doneCh := make(chan bool)
	runTasks()
	go func() {
		sig := waitForEndSignal()
		logger.Log.WithField("Signal", sig).Info("Begin closing")
		doneCh <- closeResources()
	}()
	if ok := <-doneCh; ok {
		logger.Log.Info("Gracefully shutting down")
	} else {
		logger.Log.Warn("Ungraceful shutdown")
		os.Exit(-1)
	}
}

// ExitApp gracefully ends your application and does necessary cleanup afterwards
func ExitApp() {
	sendEndSignal()
}
