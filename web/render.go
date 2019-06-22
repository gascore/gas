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

		err = replaceChild(_parent, _new, _old, node.New, node.Old)
		if err != nil {
			return err
		}
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

		err = appendChild(_parent, _new, node.New)
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

		err := removeChild(_parent, _old, node.Old)
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

		e.RChildes = []interface{}{}

		return e.Update()
	}

	return nil
}

// ChildNodes return *dom.Element child nodes
func (w BackEnd) ChildNodes(node interface{}) []interface{} {
	var iNodes []interface{}

	_childes := node.(*dom.Element).ChildNodes()
	for _, _el := range _childes {
		iNodes = append(iNodes, _el)
	}

	return iNodes
}

func replaceChild(_p, _new, _old dom.Node, new interface{}, old interface{}) error {
	err := gas.CallBeforeDestroyIfCan(old)
	if err != nil {
		return err
	}

	err = gas.CallBeforeUpdateIfCan(new)
	if err != nil {
		return err
	}

	err = removeFromRefs(old)
	if err != nil {
		return err
	}

	_p.ReplaceChild(_new, _old)

	err = gas.CallUpdatedIfCan(new)
	if err != nil {
		return err
	}

	err = gas.CallMountedIfCan(new)
	if err != nil {
		return err
	}

	return nil
}

func appendChild(_p, _c dom.Node, new interface{}) error {
	err := gas.CallBeforeUpdateIfCan(new)
	if err != nil {
		return err
	}

	_p.AppendChild(_c)

	err = gas.CallUpdatedIfCan(new)
	if err != nil {
		return err
	}

	err = gas.CallMountedIfCan(new)
	if err != nil {
		return err
	}

	return nil
}

func removeChild(_p, _e dom.Node, old interface{}) error {
	err := gas.CallBeforeUpdateIfCan(old)
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

	err = gas.CallUpdatedIfCan(old)
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
