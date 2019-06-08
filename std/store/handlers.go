package store

// Handlers type for multi-files handlers storing
type Handlers map[string]Handler

// NewHandlers create new empty handlers map
func NewHandlers() Handlers {
	return make(map[string]Handler)
}

// Add add one handler to handlers
func (h Handlers) Add(name string, handler Handler) {
	if handler == nil {
		return
	}

	h[name] = handler
}

// AddAdd add many handler to handlers
func (h Handlers) AddMany(handlers map[string]Handler) {
	for key, value := range handlers {
		h[key] = value
	}
}
