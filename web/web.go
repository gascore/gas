package web

import (
	"errors"
	"fmt"

	"github.com/gascore/dom"
	"github.com/gascore/gas"
)

// BackEnd backend for core library
type BackEnd struct {
	newRenderer bool
	queue chan *gas.RenderTask
}

// Init initialize gas application
func Init(c *gas.Component, startPoint string) error {
	return InitCustom(c,startPoint,false)
}

func InitCustom(c *gas.C, startPoint string, newRenderer bool) error {
	_startEl := dom.Doc.GetElementById(startPoint)
	if _startEl == nil {
		return errors.New("invalid startPoint")
	}

	be := &BackEnd{queue: make(chan *gas.RenderTask),newRenderer:newRenderer}
	if newRenderer {
		go be.Executor()
	}

	gas := gas.New(c, be)
	_startEl.SetAttribute("data-i", gas.UUID)
	gas.Update()

	return nil
}

// GetElement get dom.Element by element
func (w *BackEnd) GetElement(c *gas.Element) interface{} {
	return dom.Doc.QuerySelector(fmt.Sprintf("[data-i='%s']", c.UUID))
}

// ConsoleLog console.log(a)
func (w *BackEnd) ConsoleLog(a ...interface{}) { dom.ConsoleLog(a...) }

// ConsoleError console.error(a)
func (w *BackEnd) ConsoleError(a ...interface{}) { dom.ConsoleError(a...) }
