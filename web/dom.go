package web

import (
	"errors"
	"fmt"

	"github.com/gascore/dom"
	"github.com/gascore/gas"
)

// BackEnd backend for core library
type BackEnd struct{}

// CanRender check if root element exists
func (w BackEnd) CanRender(startPoint string) (string, error) {
	_el := dom.Doc.GetElementById(startPoint)
	if _el == nil {
		return "", errors.New("invalid startPoint")
	}

	return _el.TagName(), nil
}

// Init initialize gas app
func (w BackEnd) Init(gas gas.Gas) error {
	dom.Doc.GetElementById(gas.StartPoint).SetAttribute("data-i", gas.App.UUID)

	gas.App.Update()

	return nil
}

// GetElement get dom.Element by element
func (w BackEnd) GetElement(c *gas.Element) interface{} {
	return dom.Doc.QuerySelector(fmt.Sprintf("[data-i='%s']", c.UUID))
}

// GetGasEl get root dom.Element
func (w BackEnd) GetGasEl(g *gas.Gas) interface{} {
	_gas := dom.Doc.GetElementById(g.StartPoint)
	if _gas == nil {
		dom.ConsoleError("GetGasEl returning nil")
		return nil
	}

	return _gas
}

// GetBackEnd return Backend
func GetBackEnd() gas.BackEnd {
	return BackEnd{}
}

// ConsoleLog console.log(a)
func (w BackEnd) ConsoleLog(a ...interface{}) { dom.ConsoleLog(a...) }

// ConsoleError console.error(a)
func (w BackEnd) ConsoleError(a ...interface{}) { dom.ConsoleError(a...) }
