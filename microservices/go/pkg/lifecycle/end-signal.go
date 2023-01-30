package lifecycle

import (
	"os"
	"os/signal"
	"syscall"
)

var _endSignalCh = make(chan os.Signal, 1)

func init() {
	signal.Notify(_endSignalCh, syscall.SIGINT, syscall.SIGTERM)
}

func sendEndSignal() {
	_endSignalCh <- syscall.SIGINT
}

func waitForEndSignal() os.Signal {
	return <-_endSignalCh
}
