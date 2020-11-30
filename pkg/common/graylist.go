package common

import (
	"sync"
	"time"
)

// Graylist .
type Graylist struct {
	states sync.Map
}

// IsAccessable .
func (ts *Graylist) IsAccessable(name string) bool {
	value, exists := ts.states.Load(name)

	return !exists || value.(int64)-time.Now().Unix() < 0
}

// Lock .
func (ts *Graylist) Lock(name string, duration time.Duration) {
	ts.states.Store(name, time.Now().Add(duration).Unix())
}
