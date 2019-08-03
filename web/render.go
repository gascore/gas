package web

import (
	"errors"

	"github.com/gascore/dom"
	"github.com/gascore/gas"
)

// ExecNode execute render node
func (w BackEnd) ExecNode(node *gas.RenderNode) error {
	hook := func(f func(interface{}) error, el interface{}) {
		if !node.IgnoreHooks {
			err := f(el)
			if err != nil {
				dom.ConsoleError(err.Error())
			}
		}
	}

	if !node.IgnoreHooks {
		err := gas.CallBeforeUpdate(node.Parent)
		if err != nil {
			dom.ConsoleError(err.Error())
		}
	}

	switch node.Type {
	case gas.ReplaceType:
		hook(gas.CallBeforeCreated, node.New)
		hook(gas.CallBeforeDestroy, node.Old)

		// custom logic
		if node.ReplaceCanGoDeeper { // Update old element attributes
			newE := node.New.(*gas.E)
			_old := node.NodeOld.(*dom.Element)

			for attrKey, attrVal := range gas.DiffAttrs(newE.RAttrs, node.Old.(*gas.E).RAttrs) {
				_old.SetAttribute(attrKey, attrVal)
				if attrKey == "value" {
					_old.SetValue(attrVal)
				}
			}

			_old.SetAttribute("data-i", newE.UUID)
			return nil
		}

		// Create new element and replace with old
		_new, err := CreateElement(node.New)
		if err != nil {
			return err
		}

		_parent, ok := node.NodeParent.(*dom.Element)
		if !ok {
			return errors.New("invalid NodeParent type")
		}

		_old, ok := node.NodeOld.(*dom.Element)
		if !ok {
			return errors.New("invalid NodeOld type")
		}

		_parent.ReplaceChild(_new, _old)

		hook(gas.CallMounted, node.New)
	case gas.ReplaceHooks:
		hook(gas.CallMounted, node.New)
	case gas.CreateType:
		_parent, ok := node.NodeParent.(*dom.Element)
		if !ok {
			return errors.New("invalid NodeParent type")
		}

		hook(gas.CallBeforeCreated, node.New)

		_new, err := CreateElement(node.New)
		if err != nil {
			return err
		}

		_parent.AppendChild(_new)

		hook(gas.CallMounted, node.New)
	case gas.DeleteType:
		_parent, ok := node.NodeParent.(*dom.Element)
		if !ok {
			return errors.New("invalid NodeParent")
		}

		_old, ok := node.NodeOld.(*dom.Element)
		if !ok {
			return errors.New("invalid NodeOld")
		}

		hook(gas.CallBeforeDestroy, node.Old)

		_parent.RemoveChild(_old)
	case gas.RecreateType:
		e := node.New.(*gas.Element)
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

	if !node.IgnoreHooks {
		err := gas.CallUpdated(node.Parent)
		if err != nil {
			dom.ConsoleError(err.Error())
		}
	}

	return nil
}

// ChildNodes return *dom.Element child nodes
func (w BackEnd) ChildNodes(node interface{}) []interface{} {
	var iNodes []interface{}

	_node, ok := node.(*dom.Element)
	if !ok {
		return iNodes
	}

	_childes := _node.ChildNodes()
	for _, _el := range _childes {
		iNodes = append(iNodes, _el)
	}

	return iNodes
}
