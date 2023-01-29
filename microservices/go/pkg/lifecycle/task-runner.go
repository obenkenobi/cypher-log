package lifecycle

import "sync"

// TaskRunner is a general interface to run a task
type TaskRunner interface {
	// Run runs the given task
	Run()
}

var taskRunners []TaskRunner
var taskRunnerMutex sync.Mutex

// RegisterTaskRunner registers a TaskRunner in your app lifecycle to be run
// concurrently
func RegisterTaskRunner(taskRunner TaskRunner) {
	taskRunnerMutex.Lock()
	taskRunnerMutex.Unlock()
	taskRunners = append(taskRunners, taskRunner)
}

func runTasks() {
	for _, taskRunner := range taskRunners {
		go taskRunner.Run()
	}
}
