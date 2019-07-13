package gas

// Gas - main application struct
type Gas struct {
	App        *Element
	StartPoint string // html element id where application will store
}

// GetElement return root element
func (g *Gas) GetElement() interface{} {
	return g.App.RC.BE.GetGasEl(g)
}

// BackEnd interface for calling platform-specific code
type BackEnd interface {
	CanRender(string) (string, error)
	Init(Gas) error

	ExecNode(*RenderNode) error
	ChildNodes(interface{}) []interface{}

	GetElement(*Element) interface{}
	GetGasEl(*Gas) interface{}

	EditWatcherValue(interface{}, string)

	ConsoleLog(...interface{})
	ConsoleError(...interface{})
}
