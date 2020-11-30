package common

import (
	s "sync"
	"testing"
	"time"
)

func TestSemaphore(t *testing.T) {
	wg := s.WaitGroup{}
	sem := NewSemaphore(5)

	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			sem.Enter(1)
			defer wg.Done()
			defer sem.Leave(1)

			time.Sleep(100 * time.Millisecond)
		}()
	}

	startTime := time.Now()

	wg.Wait()

	diff := (time.Now().UnixNano() - startTime.UnixNano()) / 1000000
	if diff < 150 || diff > 250 {
		t.Fatalf("Expected 150~250 msec, got %v", diff)
	}
}
