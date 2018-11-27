// Package gas -- main package
package gas

import (
	"errors"

	dom "github.com/dennwc/dom"
)

// Gas -- struct for main application component
type Gas struct {
	App        Component
	StartPoint string // html element id where application will store
	Element    *dom.Element

	// Other stuff
}

// New create new gas application
//
//
// All Add* methods return current component, if they don't we need to pre create all components before run NewComponent
// it would looks like:
//
// `
//	c1c1 := NewComponent(...)
//	c1c2 := NewComponent(...)
//
// 	c1 := NewComponent(...).AddChildes(c1c1, c1c2)
//	c2 := NewComponent(...)
//	c3 := NewComponent(...)
//
//	component := NewComponent(...).Add*(...).AddChildes(c1, c2, c3)
// ` -- seems little ridiculous
func New(startPoint string, components ...GetComponent) (Gas, error) {
	el := dom.GetDocument().GetElementById(startPoint)
	if el == nil {
		return Gas{}, errors.New("invalid start point")
	}

	mainComponent := NewComponent(NilData, NilData, "wrap", startPoint, NilClasses, NilAttrs).AddBinds(NilBinds).AddChildes(components...)

	gas := Gas{App: *mainComponent, StartPoint: startPoint, Element: el}

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

//func must(err error) {
//	if err != nil {
//		panic(err)
//	}
//}
