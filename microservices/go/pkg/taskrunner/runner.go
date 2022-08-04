package taskrunner

// TaskRunner is a general interface to run a task
type TaskRunner interface {
	// Run runs the given task
	Run()
}
