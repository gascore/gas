package gas

import (
	"errors"

	"github.com/Sinicablyat/dom"
)

// Gas -- struct for main application component
type Gas struct {
	App        Component
	StartPoint string // html element id where application will store
	Element    *dom.Element

	// Other stuff
}

// New create new gas application
func New(startPoint string, components ...GetComponent) (Gas, error) {
	_el := dom.Doc.GetElementById(startPoint)
	if _el == nil {
		return Gas{}, errors.New("invalid start point")
	}

	mainComponent := NewComponent(
		NilParentComponent,
		NilData, // data
		NilData, // props
		NilMethods, // methods
		NilBinds, // binds
		NilHandlers, // handlers
		_el.GetTagName(), // tag name
		map[string]string{ // attributes
			"id": startPoint,
		},
		components...) // components

	gas := Gas{App: *mainComponent, StartPoint: startPoint, Element: _el}

	return gas, nil
}

// Init initialize gas application
func (gas *Gas) Init() error {
	app := gas.App
	_main := gas.Element

	for _, el := range app.Childes(app) {
		_child, err := CreateComponent(el)
		if err != nil {
			return err
		}

		_main.AppendChild(_child)
	}

	return nil
}

var signal = make(chan int)
func KeepAlive() {
	for {
		<-signal
	}
}

//func must(err error) {
//	if err != nil {
//		panic(err)
//	}
//}
