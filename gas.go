package gas

import "fmt"

// Gas - main application struct
type Gas struct {
	Component *Component
}

// BackEnd interface for calling platform-specific code
type BackEnd interface {
	ExecTasks([]*RenderTask)
	GetElement(*Element) interface{}
	ChildNodes(interface{}) []interface{}

	ConsoleLog(...interface{})
	ConsoleError(...interface{})
}

// New create new gas application with custom backend
func New(c *Component, be BackEnd) *Gas {
	c.RC = &RenderCore{BE: be}
	return &Gas{c}
}

// ToGetComponentList return array by many parameters, because it's pretty
func ToGetComponentList(childes ...interface{}) []interface{} {
	return childes
}

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
