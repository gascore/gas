package gas

import "fmt"

// New create new gas application with custom backend
func New(c *Component, be BackEnd) Gas {
	c.RC = &RenderCore{BE: be}
	return Gas{c.Init()}
}

// ToGetComponentList return array by many parameters, because it's pretty
func ToGetComponentList(childes ...interface{}) []interface{} {
	return childes
}

// CL alias for ToGetComponentList
var CL = ToGetComponentList

// WarnError log error
func (c *Component) WarnError(err error) {
	if err == nil {
		return
	}

	c.ConsoleError(err.Error())
}

// WarnIfNot console error if !ok
func (c *Component) WarnIfNot(ok bool) {
	if ok {
		return
	}

	c.ConsoleError(fmt.Errorf("invalid Data type").Error())
}

// ConsoleLog call BackEnd.ConsoleLog
func (c *Component) ConsoleLog(a ...interface{}) { c.RC.BE.ConsoleLog(a...) }

// ConsoleError call BackEnd.ConsoleError
func (c *Component) ConsoleError(a ...interface{}) { c.RC.BE.ConsoleError(a...) }
