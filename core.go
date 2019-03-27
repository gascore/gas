package gas

// Gas - main application struct
type Gas struct {
	App        *Component
	StartPoint string // html element id where application will store

	// Other stuff
}

// GetElement return root element
func (g *Gas) GetElement() interface{} {
	return g.App.RC.BE.GetGasEl(g)
}

// BackEnd interface for calling platform-specific code
type BackEnd interface {
	New(string) (string, error)
	Init(Gas) error

	ExecNode(*RenderNode) error
	ChildNodes(interface{}) []interface{}

	GetElement(*Component) interface{}
	GetGasEl(*Gas) interface{}

	ConsoleLog(...interface{})
	ConsoleError(...interface{})
}
