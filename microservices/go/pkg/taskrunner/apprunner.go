package taskrunner

import "sync"

// RunAndWait runs each tasks concurrently and waits for each of them to complete
func RunAndWait(taskRunners ...TaskRunner) {
	taskRunnersLength := len(taskRunners)
	if taskRunnersLength == 0 {
		return
	}
	var wg sync.WaitGroup
	for _, taskRunner := range taskRunners[1:taskRunnersLength] {
		taskRunner := taskRunner
		wg.Add(1)
		go func() {
			defer wg.Done()
			taskRunner.Run()
		}()
	}
	taskRunners[0].Run()
	wg.Wait()
}
