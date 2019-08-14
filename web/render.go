package web

import (
	"errors"

	"github.com/gascore/dom"
	"github.com/gascore/gas"
)

func (w *BackEnd) Executor() {
	for {
		select {
		case task := <- w.queue:
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
		if !task.IgnoreHooks {
			err := f(el)
			if err != nil {
				dom.ConsoleError(err.Error())
			}
		}
	}

	if !task.IgnoreHooks {
		err := gas.CallBeforeUpdate(task.Parent)
		if err != nil {
			dom.ConsoleError(err.Error())
		}
	}

	switch task.Type {
	case gas.ReplaceType:
		hook(gas.CallBeforeCreated, task.New)
		hook(gas.CallBeforeDestroy, task.Old)

		// custom logic
		if task.ReplaceCanGoDeeper { // Update old element attributes
			newE := task.New.(*gas.E)
			_old := task.NodeOld.(*dom.Element)

			for attrKey, attrVal := range gas.DiffAttrs(newE.RAttrs, task.Old.(*gas.E).RAttrs) {
				_old.SetAttribute(attrKey, attrVal)
				if attrKey == "value" {
					_old.SetValue(attrVal)
				}
			}

			_old.SetAttribute("data-i", newE.UUID)
			return nil
		}

		// Create new element and replace with old
		_new, err := CreateElement(task.New)
		if err != nil {
			return err
		}

		_parent, ok := task.NodeParent.(*dom.Element)
		if !ok {
			return errors.New("invalid NodeParent type")
		}

		_old, ok := task.NodeOld.(*dom.Element)
		if !ok {
			return errors.New("invalid NodeOld type")
		}

		_parent.ReplaceChild(_new, _old)

		hook(gas.CallMounted, task.New)
	case gas.ReplaceHooks:
		hook(gas.CallMounted, task.New)
	case gas.CreateType:
		_parent, ok := task.NodeParent.(*dom.Element)
		if !ok {
			return errors.New("invalid NodeParent type")
		}

		hook(gas.CallBeforeCreated, task.New)

		_new, err := CreateElement(task.New)
		if err != nil {
			return err
		}

		_parent.AppendChild(_new)

		hook(gas.CallMounted, task.New)
	case gas.DeleteType:
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
	case gas.RecreateType:
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

		err := e.Update()
		if err != nil {
			return err
		}
	}

	if !task.IgnoreHooks {
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
