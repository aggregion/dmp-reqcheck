package common

import "context"

type (
	empty struct{}

	// Semaphore for concurrency
	Semaphore chan empty
)

var (
	emptyValue empty
)

// NewSemaphore create semaphore
func NewSemaphore(size int) Semaphore {
	return make(Semaphore, size)
}

// Enter inc owners
func (smp *Semaphore) Enter(n int) {
	for i := 0; i < n; i++ {
		*smp <- emptyValue
	}
}

// EnterWithCtx inc owners with ctx
func (smp *Semaphore) EnterWithCtx(ctx context.Context, n int) {
	for i := 0; i < n; i++ {
		select {
		case <-ctx.Done():
			return
		case *smp <- emptyValue:
		}
	}
}

// Leave dec owners
func (smp *Semaphore) Leave(n int) {
	for i := 0; i < n; i++ {
		<-*smp
	}
}

// LeaveWithCtx dec owners with ctx
func (smp *Semaphore) LeaveWithCtx(ctx context.Context, n int) {
	for i := 0; i < n; i++ {
		select {
		case <-ctx.Done():
			return
		case <-*smp:
		}
	}
}
