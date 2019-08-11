package gas

// Gas - main application struct
type Gas struct {
	*Element
}

// BackEnd interface for calling platform-specific code
type BackEnd interface {
	ExecNode(*RenderNode) error
	GetElement(*Element) interface{}
	ChildNodes(interface{}) []interface{}

	ConsoleLog(...interface{})
	ConsoleError(...interface{})
}
