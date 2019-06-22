package gas

import (
	"fmt"
	"sync"
	"github.com/eapache/queue"
)

// New create new gas application with custom backend
func New(be BackEnd, startPoint string, c *Component) (Gas, error) {
	q := queue.New()
	c.RC = &RenderCore{
		BE:    be,
		WG:    &sync.WaitGroup{},
		Queue: q,
	}

	_, err := be.CanRender(startPoint)
	if err != nil {
		return Gas{}, err
	}

	gas := Gas{App: c.Init(), StartPoint: startPoint}

	return gas, nil
}

// Init initialize gas application
func Init(be BackEnd, startPoint string, c *Component) error {
	gas, err := New(be, startPoint, c)
	if err != nil {
		return err
	}

	err = be.Init(gas)
	if err != nil {
		return err
	}

	return nil
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
