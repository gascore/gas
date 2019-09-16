package web

import (
	"errors"
	// "fmt"

	"github.com/gascore/dom"
	"github.com/gascore/gas"
)

func (w *BackEnd) Executor() {
	for {
		select {
		case task := <-w.queue:
			err := w.ExecTask(task)
			if err != nil {
				dom.ConsoleError(err.Error())
			}
		}
	}
}

// ExecTasks execute render tasks
func (w *BackEnd) ExecTasks(tasks []*gas.RenderTask) {
	if w.newRenderer {
		for _, task := range tasks {
			w.queue <- task
		}
		return
	}

	for _, task := range tasks {
		err := w.ExecTask(task)
		if err != nil {
			dom.ConsoleError(err.Error())
		}
	}
}

// ExecTask execute render task
func (w *BackEnd) ExecTask(task *gas.RenderTask) error {
	hook := func(f func(interface{}) error, el interface{}) {
		if !task.InReplaced {
			err := f(el)
			if err != nil {
				dom.ConsoleError(err.Error())
			}
		}
	}

	if !task.InReplaced {
		err := gas.CallBeforeUpdate(task.Parent)
		if err != nil {
			dom.ConsoleError(err.Error())
		}
	}

	switch task.Type {
	case gas.RReplace:
		hook(gas.CallBeforeCreated, task.New)
		hook(gas.CallBeforeDestroy, task.Old)

		// Update old element attributes
		if task.ReplaceCanGoDeeper {
			newE := task.New.(*gas.E)
			_old := task.NodeOld.(*dom.Element)

			setAttributes(_old, gas.DiffAttrs(newE.RAttrs, task.Old.(*gas.E).RAttrs))
			_old.SetAttribute("data-i", newE.UUID)

			return nil
		}

		// Create new element and replace with old
		_new, err := CreateElement(task.New)
		if err != nil {
			return err
		}

		_old, ok := task.NodeOld.(*dom.Element)
		if !ok {
			return errors.New("invalid NodeOld type")
		}

		_old.ParentElement().ReplaceChild(_new, _old)

		hook(gas.CallMounted, task.New)
	case gas.RReplaceHooks:
		hook(gas.CallMounted, task.New)
	case gas.RCreate, gas.RFirstRender:
		var _parent *dom.Element
		if task.Type == gas.RCreate {
			var ok bool
			_parent, ok = task.NodeParent.(*dom.Element)
			if !ok {
				return errors.New("invalid NodeParent type")
			}
		} else {
			_parent = w.startEl
		}

		hook(gas.CallBeforeCreated, task.New)

		_new, err := CreateElement(task.New)
		if err != nil {
			return err
		}

		_parent.AppendChild(_new)

		hook(gas.CallMounted, task.New)
	case gas.RDelete:
		_parent, ok := task.NodeParent.(*dom.Element)
		if !ok {
			return errors.New("invalid NodeParent")
		}

		_old, ok := task.NodeOld.(*dom.Element)
		if !ok {
			return errors.New("invalid NodeOld")
		}

		hook(gas.CallBeforeDestroy, task.Old)

		_parent.RemoveChild(_old)
	case gas.RRecreate:
		e := task.New.(*gas.Element)
		_e, ok := e.BEElement().(*dom.Element)
		if !ok {
			return errors.New("invalid NodeNew type")
		}

		hook(gas.CallBeforeDestroy, e)

		for _, _child := range _e.ChildNodes() {
			_e.RemoveChild(_child)
		}

		e.Childes = []interface{}{}
		e.OldChildes = []interface{}{}

		err := e.ParentComponent().Component.UpdateWithError()
		if err != nil {
			return err
		}
	}

	if !task.InReplaced {
		err := gas.CallUpdated(task.Parent)
		if err != nil {
			dom.ConsoleError(err.Error())
		}
	}

	return nil
}

// ChildNodes return *dom.Element child nodes
func (w *BackEnd) ChildNodes(el interface{}) []interface{} {
	var iChildes []interface{}
	_el, ok := el.(*dom.Element)
	if !ok {
		return iChildes
	}

	_childes := _el.ChildNodes()
	for _, _el := range _childes {
		iChildes = append(iChildes, _el)
	}

	return iChildes
}
