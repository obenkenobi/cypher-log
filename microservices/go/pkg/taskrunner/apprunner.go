package taskrunner

import "sync"

// RunAndWait runs each tasks in a goroutine and waits for each of them to complete
func RunAndWait(tasks ...func()) {
	var wg sync.WaitGroup
	for _, task := range tasks {
		task := task
		wg.Add(1)
		go func() {
			defer wg.Done()
			task()
		}()
	}
	wg.Wait()
}
