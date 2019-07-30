package web

import (
	"errors"

	"github.com/gascore/dom"
	"github.com/gascore/gas"
)

// ExecNode execute render node
func (w BackEnd) ExecNode(node *gas.RenderNode) error {
	switch node.Type {
	case gas.ReplaceType:
		if differOnlyInAttributes(node.New, node.Old) {
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

		err := gas.CallBeforeCreatedIfCan(node.New)
		if err != nil {
			return nil
		}

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

		err = replaceChild(_parent, _new, _old, node.New, node.Old, node.Parent)
		if err != nil {
			return err
		}

		return nil
	case gas.CreateType:
		err := gas.CallBeforeCreatedIfCan(node.New)
		if err != nil {
			return nil
		}

		_new, err := CreateElement(node.New)
		if err != nil {
			return err
		}

		_parent, ok := node.NodeParent.(*dom.Element)
		if !ok {
			return errors.New("invalid NodeParent type")
		}

		err = appendChild(_parent, _new, node.New, node.Parent)
		if err != nil {
			return err
		}

		return nil
	case gas.DeleteType:
		_parent, ok := node.NodeParent.(*dom.Element)
		if !ok {
			return errors.New("invalid NodeParent")
		}

		_old, ok := node.NodeOld.(*dom.Element)
		if !ok {
			return errors.New("invalid NodeOld")
		}

		err := removeChild(_parent, _old, node.Old, node.Parent)
		if err != nil {
			return err
		}

		return nil
	case gas.RecreateType:
		e := node.New.(*gas.Element)
		_e, ok := e.BEElement().(*dom.Element)
		if !ok {
			return errors.New("invalid NodeNew type")
		}

		err := gas.CallBeforeDestroyIfCan(e)
		if err != nil {
			return err
		}

		for _, _child := range _e.ChildNodes() {
			_e.RemoveChild(_child)
		}

		e.Childes = []interface{}{}
		e.OldChildes = []interface{}{}

		return e.Update()
	}

	return nil
}

// differOnlyInAttributes return true if only defference between elements is Attrs
func differOnlyInAttributes(new, old interface{}) bool {
	newE, ok := new.(*gas.Element)
	if !ok {
		return false
	}

	oldE, ok := old.(*gas.Element)
	if !ok {
		return false
	}

	if len(newE.Childes) != 0 || len(newE.OldChildes) != 0 {
		return false
	}

	return gas.ElementsCanBeUpdated(newE, oldE)
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

func replaceChild(_p, _new, _old dom.Node, new interface{}, old interface{}, p *gas.E) error {
	err := gas.CallBeforeDestroyIfCan(old)
	if err != nil {
		return err
	}

	err = gas.CallBeforeUpdateIfCan(new, p)
	if err != nil {
		return err
	}

	err = removeFromRefs(old)
	if err != nil {
		return err
	}

	_p.ReplaceChild(_new, _old)

	err = gas.CallUpdatedIfCan(new, p)
	if err != nil {
		return err
	}

	err = gas.CallMountedIfCan(new)
	if err != nil {
		return err
	}

	return nil
}

func appendChild(_p, _c dom.Node, new interface{}, p *gas.E) error {
	err := gas.CallBeforeUpdateIfCan(new, p)
	if err != nil {
		return err
	}

	_p.AppendChild(_c)

	err = gas.CallUpdatedIfCan(new, p)
	if err != nil {
		return err
	}

	err = gas.CallMountedIfCan(new)
	if err != nil {
		return err
	}

	return nil
}

func removeChild(_p, _e dom.Node, old interface{}, p *gas.E) error {
	err := gas.CallBeforeUpdateIfCan(old, p)
	if err != nil {
		return err
	}

	err = gas.CallBeforeDestroyIfCan(old)
	if err != nil {
		return err
	}

	err = removeFromRefs(old)
	if err != nil {
		return err
	}

	_p.RemoveChild(_e)

	err = gas.CallUpdatedIfCan(old, p)
	if err != nil {
		return err
	}

	return nil
}

func removeFromRefs(old interface{}) error {
	if gas.IsElement(old) {
		oldE := gas.I2E(old)
		if len(oldE.RefName) != 0 {
			p := oldE.ParentComponent()
			if p.Component.Refs == nil {
				return errors.New("cannot remove element from parent refs because parent refs is nil")
			}
			if _, ok := p.Component.Refs[oldE.RefName]; ok {
				delete(p.Component.Refs, oldE.RefName)
			}
		}
	}
	return nil
}
