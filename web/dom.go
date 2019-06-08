package web

import (
	"errors"
	"fmt"

	"github.com/gascore/dom"
	"github.com/gascore/gas"
)

var (
	// ErrCookieNotFound if cookie not found
	ErrCookieNotFound = errors.New("cookie not found")
	// ErrInvalidCookie if cookie is invalid
	ErrInvalidCookie = errors.New("invalid cookie")
)

// BackEnd backend for core library
type BackEnd struct{}

// New check if root element exists
func (w BackEnd) New(startPoint string) (string, error) {
	_el := dom.Doc.GetElementById(startPoint)
	if _el == nil {
		return "", errors.New("invalid startPoint")
	}

	return _el.TagName(), nil
}

// Init initialize gas app
func (w BackEnd) Init(gas gas.Gas) error {
	app := gas.App
	_main := dom.Doc.GetElementById(gas.StartPoint)

	_main.SetAttribute("data-i", app.UUID)

	err := app.ForceUpdate()
	if err != nil {
		return err
	}

	dom.Doc.GetElementsByTagName("body")[0].SetAttribute("data-ready", true)

	return nil
}

// GetElement get dom.Element by component
func (w BackEnd) GetElement(c *gas.Component) interface{} {
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
