package lifecycle

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

// RunApp starts your application concurrently by running each task in
// taskRunners concurrently and closes registered resources when the app ends
func RunApp() {
	doneCh := make(chan bool)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	runTasks()

	go func() {
		sig := <-sigCh
		logger.Log.WithField("Signal", sig).Info("Begin closing")
		doneCh <- closeResources()
	}()
	if ok := <-doneCh; ok {
		logger.Log.Info("Gracefully shutting down")
	} else {
		logger.Log.Warn("Ungraceful shutdown")
	}
}
