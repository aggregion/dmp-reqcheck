package utils

type (
	// EventHandler .
	EventHandler = func(args ...interface{}) error

	// Event .
	Event interface {
		On(handler EventHandler)
		Emit(args ...interface{}) error
	}

	// eventImpl .
	eventImpl struct {
		handlers []EventHandler
	}
)

// Emit .
func (ev *eventImpl) Emit(args ...interface{}) error {
	var err error
	for _, handler := range ev.handlers {
		err = handler(args...)
		if err != nil {
			return err
		}
	}

	return nil
}

// On .
func (ev *eventImpl) On(handler EventHandler) {
	ev.handlers = append(ev.handlers, handler)
}

// NewSyncEvent creates default implementation
func NewSyncEvent() Event {
	return &eventImpl{
		make([]EventHandler, 0, 1),
	}
}
