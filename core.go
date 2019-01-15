package gas

// Gas - main application struct
type Gas struct {
	App        *Component
	StartPoint string // html element id where application will store

	// Other stuff
}

func (g *Gas) GetElement() interface{} {
	return g.App.be.GetGasEl(g)
}

type BackEnd interface {
	New(string) (string, error)
	Init(Gas) error
	UpdateComponentChildes(*Component, []interface{}, []interface{}) error
	ReCreate(*Component) error
	RenderTree(*Component) []interface{}

	GetElement(*Component) interface{}
	GetGasEl(*Gas) interface{}

	ConsoleLog(...interface{})
	ConsoleError(...interface{})
}