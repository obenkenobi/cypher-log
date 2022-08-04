package taskrunner

import "sync"

// RunAndWait runs each tasks in a goroutine and waits for each of them to complete
func RunAndWait(taskRunners ...TaskRunner) {
	var wg sync.WaitGroup
	for _, taskRunner := range taskRunners {
		taskRunner := taskRunner
		wg.Add(1)
		go func() {
			defer wg.Done()
			taskRunner.Run()
		}()
	}
	wg.Wait()
}
