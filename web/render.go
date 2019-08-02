package web

import (
	"errors"

	"github.com/gascore/dom"
	"github.com/gascore/gas"
)

// ExecNode execute render node
func (w BackEnd) ExecNode(node *gas.RenderNode) error {
	err := gas.CallBeforeUpdate(node.Parent)
	if err != nil {
		return err
	}

	switch node.Type {
	case gas.ReplaceType:
		err := gas.CallBeforeCreated(node.New)
		if err != nil {
			return nil
		}

		err = gas.CallBeforeDestroy(node.Old)
		if err != nil {
			return err
		}

		// custom logic
		if node.ReplaceCanGoDeeper {
			// Update old element attributes
			newE := node.New.(*gas.E)
			_old := node.NodeOld.(*dom.Element)

			for attrKey, attrVal := range gas.DiffAttrs(newE.RAttrs, node.Old.(*gas.E).RAttrs) {
				_old.SetAttribute(attrKey, attrVal)
				if attrKey == "value" {
					_old.SetValue(attrVal)
				}
			}

			_old.SetAttribute("data-i", newE.UUID)
		} else {
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
		}

		err = gas.CallMounted(node.New)
		if err != nil {
			return err
		}
	case gas.CreateType:
		_parent, ok := node.NodeParent.(*dom.Element)
		if !ok {
			return errors.New("invalid NodeParent type")
		}

		err := gas.CallBeforeCreated(node.New)
		if err != nil {
			return nil
		}

		_new, err := CreateElement(node.New)
		if err != nil {
			return err
		}

		_parent.AppendChild(_new)

		err = gas.CallMounted(node.New)
		if err != nil {
			return err
		}
	case gas.DeleteType:
		_parent, ok := node.NodeParent.(*dom.Element)
		if !ok {
			return errors.New("invalid NodeParent")
		}

		_old, ok := node.NodeOld.(*dom.Element)
		if !ok {
			return errors.New("invalid NodeOld")
		}

		err = gas.CallBeforeDestroy(node.Old)
		if err != nil {
			return err
		}

		_parent.RemoveChild(_old)
	case gas.RecreateType:
		e := node.New.(*gas.Element)
		_e, ok := e.BEElement().(*dom.Element)
		if !ok {
			return errors.New("invalid NodeNew type")
		}

		err := gas.CallBeforeDestroy(e)
		if err != nil {
			return err
		}

		for _, _child := range _e.ChildNodes() {
			_e.RemoveChild(_child)
		}

		e.Childes = []interface{}{}
		e.OldChildes = []interface{}{}

		err = e.Update()
		if err != nil {
			return err
		}
	}

	err = gas.CallUpdated(node.Parent)
	if err != nil {
		return err
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
