package gas

// emptyBackEnd empty backend for testing backend calling methods
type emptyBackEnd struct{}

func GetEmptyRenderCore() *RenderCore {
	return &RenderCore{BE:emptyBackEnd{}}
}

func GetEmptyBackend() BackEnd {
	return emptyBackEnd{}
}

func (e emptyBackEnd) ExecNode(node *RenderNode) error {
	return nil
}

func (e emptyBackEnd) ChildNodes(i interface{}) []interface{} {
	return []interface{}{}
}

func (e emptyBackEnd) New(a string) (string, error) {
	return "app", nil
}

func (e emptyBackEnd) Init(g Gas) error {
	return nil
}

func (e emptyBackEnd) UpdateComponentChildes(c *Component, newChildes, oldChildes []interface{}) error {
	return nil
}

func (e emptyBackEnd) ReCreate(c *Component) error {
	return nil
}

func (e emptyBackEnd) GetElement(c *Component) interface{} {
	return "not nil!"
}

func (e emptyBackEnd) GetGasEl(g *Gas) interface{} {
	return "not nil!"
}

func (e emptyBackEnd) ConsoleLog(values ...interface{}) {}

func (e emptyBackEnd) ConsoleError(values ...interface{}) {}
