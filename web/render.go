package web

import (
	"fmt"
	"reflect"

	"github.com/gascore/dom"
	"github.com/gascore/gas"
	"github.com/pkg/errors"
)

// ExecNode execute render node
func (w BackEnd) ExecNode(node *gas.RenderNode) error {
	switch node.Type {
	case gas.ReplaceType:
		err := gas.CallBeforeCreatedIfCan(node.New)
		if err != nil {
			return nil
		}

		_new, err := CreateComponent(node.New)
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

		_new, err := CreateComponent(node.New)
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
	case gas.SyncType:
		newC := node.New.(*gas.Component)
		_el, ok := node.NodeNew.(*dom.Element)
		if !ok {
			return errors.New("invalid NodeNew type")
		}

		if newC.Model.Component != nil { // update input value
			if len(newC.Model.Deep) == 0 {
				_el.SetValue(newC.Model.Component.Get(newC.Model.Data))
			} else {
				dataValue := newC.Model.Component.Get(newC.Model.Data)
				field, err := getField(reflect.ValueOf(dataValue), newC.Model.Deep)
				if err != nil {
					return err
				}

				_el.SetValue(field.Interface())
			}
		}

		updateVisible(newC, _el)
		return nil
	case gas.RecreateType:
		c := node.New.(*gas.Component)
		_c, ok := c.Element().(*dom.Element)
		if !ok {
			return errors.New("invalid NodeNew type")
		}

		err := gas.CallBeforeDestroyIfCan(c)
		if err != nil {
			return err
		}

		for _, _e := range _c.ChildNodes() {
			_c.RemoveChild(_e)
		}

		c.RChildes = []interface{}{}

		return c.ForceUpdate()
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

	if gas.IsComponent(old) {
		oldC := gas.I2C(old)
		if len(oldC.Ref) != 0 {
			p := oldC.ParentComponent()
			if p.Refs == nil {
				return errors.New("cannot remove component from parent refs because parent refs is nil")
			}
			if _, ok := p.Refs[oldC.Ref]; ok {
				delete(p.Refs, oldC.Ref)
			}
		}
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

	if gas.IsComponent(old) {
		oldC := gas.I2C(old)
		if len(oldC.Ref) != 0 {
			p := oldC.ParentComponent()
			if p.Refs == nil {
				return errors.New("cannot remove component from parent refs because parent refs is nil")
			}
			if _, ok := p.Refs[oldC.Ref]; ok {
				delete(p.Refs, oldC.Ref)
			}
		}
	}

	_p.RemoveChild(_e)

	err = gas.CallUpdatedIfCan(old)
	if err != nil {
		return err
	}

	return nil
}

func getField(data reflect.Value, arr []gas.ModelDirectiveDeepData) (reflect.Value, error) {
	var nilValue reflect.Value
	for i, deep := range arr {
		if !data.IsValid() {
			return nilValue, errors.New("data is not valid")
		}

		if data.Kind() == reflect.Ptr {
			data = data.Elem()
		}

		var field reflect.Value
		switch data.Kind() {
		case reflect.Array, reflect.Slice:
			if !deep.Brackets {
				return nilValue, fmt.Errorf("no brackets for array (%d)", i)
			}

			switch deep.Data.(type) {
			case int:
				break
			default:
				return nilValue, fmt.Errorf("error deep data type. want number, got: %T", deep.Data)
			}

			field = data.Index(deep.Data.(int))
		case reflect.Map:
			if !deep.Brackets {
				return nilValue, fmt.Errorf("no brackets for map (%d)", i)
			}
			field = data.MapIndex(reflect.ValueOf(deep.Data))
		case reflect.Struct:
			if deep.Brackets {
				return nilValue, fmt.Errorf("has brackets for struct (%d)", i)
			}

			switch deep.Data.(type) {
			case string:
				break
			default:
				return nilValue, errors.New("invalid structure field name")
			}

			field = data.FieldByName(deep.Data.(string))
		default:
			return nilValue, fmt.Errorf("invalid data type in deep %d", i)
		}

		if !field.IsValid() {
			return nilValue, errors.New("invalid field")
		}

		data = field
	}
	return data, nil
}

func updateVisible(c *gas.Component, _el *dom.Element) {
	if c.Show != nil {
		if !c.Show(c) {
			_el.Style().Set("display", "none")
		} else {
			_el.Style().Set("display", "")
		}
	}
}
