package gas

// emptyBackEnd empty backend for testing backend calling methods
type emptyBackEnd struct {
	logger func(*RenderNode)
}

// GetEmptyRenderCore return epmty render core
func GetEmptyRenderCore() *RenderCore {
	return &RenderCore{BE: emptyBackEnd{}}
}

// GetEmptyBackend return empty BackEnd
func GetEmptyBackend() BackEnd {
	return emptyBackEnd{}
}

// ExecNode return nil
func (e emptyBackEnd) ExecNode(node *RenderNode) error {
	if e.logger != nil {
		e.logger(node)
	}

	return nil
}

// ChildNodes return nil
func (e emptyBackEnd) ChildNodes(i interface{}) []interface{} {
	return []interface{}{}
}

// New return nil
func (e emptyBackEnd) CanRender(a string) (string, error) {
	return "app", nil
}

// Init return nil
func (e emptyBackEnd) Init(g Gas) error {
	return nil
}

// UpdateComponentChildes return nil
func (e emptyBackEnd) UpdateComponentChildes(c *Component, newChildes, oldChildes []interface{}) error {
	return nil
}

// ReCreate return nil
func (e emptyBackEnd) ReCreate(c *Component) error {
	return nil
}

// GetElement return not nil
func (e emptyBackEnd) GetElement(i *Element) interface{} {
	return "not nil!"
}

// GetGasEl return not nil
func (e emptyBackEnd) GetGasEl(g *Gas) interface{} {
	return "not nil!"
}

// EditWatcherValue is empty
func (e emptyBackEnd) EditWatcherValue(el interface{}, newVal string) {}

// ConsoleLog return nil
func (e emptyBackEnd) ConsoleLog(values ...interface{}) {}

// ConsoleError return nil
func (e emptyBackEnd) ConsoleError(values ...interface{}) {}
