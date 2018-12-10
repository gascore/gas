package gas

import (
	"errors"
	"fmt"

	"github.com/Sinicablyat/dom"
)

// Gas - main application struct
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
		NilMethods, // methods
		NilComputeds, // computed
		NilDirectives, // if-directive
		NilBinds, // binds
		NilHandlers, // handlers
		_el.GetTagName(), // tag name
		map[string]string{ // attributes
			"id": startPoint,
			"data-main": "true",
		},
		components...) // components

	gas := Gas{App: *mainComponent, StartPoint: startPoint, Element: _el}

	return gas, nil
}

// Init initialize gas application
func (gas *Gas) Init() error {
	app := gas.App
	_main := gas.Element

	_main.SetAttribute("data-i", app.UUID)
	for _, el := range app.Childes(&app) {
		_child, err := CreateComponent(el)
		if err != nil {
			return err
		}

		_main.AppendChild(_child)
	}

	dom.Doc.GetElementsByTagName("body")[0].SetAttribute("data-ready", true)

	return nil
}

// WarnError log error
func WarnError(err error) {
	if err == nil {
		return
	}

	dom.ConsoleError(err.Error())
}

// WarnIfNot console error if !ok
func WarnIfNot(ok bool) {
	if ok {
		return
	}

	dom.ConsoleError(fmt.Sprintf("invalid data type"))
}

var signal = make(chan int)
// KeepAlive keep alive runtime, without it application will stop (user won't be able to init events)
func KeepAlive() {
	//var wg sync.WaitGroup
	//wg.Add(1)
	//go func() {
	//	wg.Wait()
	//}()
	for {
		<- signal
	}
}

//func must(err error) {
//	if err != nil {
//		panic(err)
//	}
//}
