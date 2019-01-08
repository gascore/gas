package gas

import (
	"fmt"
)

// New create new gas application with custom backend
func New(be BackEnd, startPoint string, c *Component, getChildes GetComponentChildes) (Gas, error) {
	c.be = be

	tagName, err := c.be.New(startPoint)
	if err != nil {
		return Gas{}, err
	}

	c.Tag = tagName

	if c.Attrs == nil { c.Attrs = make(map[string]string) }
	c.Attrs["id"] = startPoint
	c.Attrs["data-main"] = "true"

	mainComponent := NewComponent(c, getChildes)

	gas := Gas{App: *mainComponent, StartPoint: startPoint}

	return gas, nil
}

// Init initialize gas application
func Init(gas Gas) error {
	err := gas.App.be.Init(gas)
	if err != nil {
		return err
	}

	return nil
}

// ToGetComponentList return array by many parameters, because it's pretty
func ToGetComponentList(childes ...interface{}) []interface{} {
	return childes
}
var CL = ToGetComponentList

// WarnError log error
func (c *Component) WarnError(err error) {
	if err == nil {
		return
	}

	c.be.ConsoleError(err.Error())
}

// WarnIfNot console error if !ok
func (c *Component) WarnIfNot(ok bool) {
	if ok {
		return
	}

	c.be.ConsoleError(fmt.Sprintf("invalid data type"))
}

func (c *Component) ConsoleLog  (a ...interface{}) { c.be.ConsoleLog  (a...) }
func (c *Component) ConsoleError(a ...interface{}) { c.be.ConsoleError(a...) }

var signal = make(chan int)

// KeepAlive keep alive runtime, without it application will stop (user won't be able to init events)
func KeepAlive() {
	for {
		<- signal
	}
}