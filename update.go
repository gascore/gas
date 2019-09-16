package gas

import (
	"errors"
	"fmt"
)

// UpdateWithError update component childes
func (component *Component) UpdateWithError() error {
	oldChild := component.Element // oldChild
	if oldChild == nil {          // it's Gas in first render
		component.RC.Add(&RenderTask{
			Type: RFirstRender,
			New:  component.RC.updateChild(nil, component.RenderElement()),
		})
		component.RC.Exec()
		return nil
	}

	_oldChild := oldChild.BEElement()
	if _oldChild == nil {
		return errors.New("invalid _oldChild")
	}

	p := oldChild.Parent

	newChild := component.RC.updateChild(p, component.RenderElement())
	isChanged, canGoDeeper, err := Changed(newChild, oldChild)
	if err != nil {
		return err
	}
	if isChanged {
		component.RC.Add(&RenderTask{
			Type:               RReplace,
			NodeOld:            _oldChild,
			New:                newChild,
			Old:                oldChild,
			ReplaceCanGoDeeper: canGoDeeper,
			Parent:             p,
		})
	}

	if canGoDeeper {
		newChildE, ok := newChild.(*Element)
		if !ok {
			return fmt.Errorf("unexpected newE type want: *gas.Element, got: %T", newChild)
		}

		newChildE.UUID = oldChild.UUID

		err := component.RC.UpdateElementChildes(_oldChild, oldChild, newChildE.Childes, getOldChildes(newChildE, oldChild), true)
		if err != nil {
			return err
		}
	}

	if isChanged {
		component.RC.Add(&RenderTask{
			Type:   RReplaceHooks,
			Parent: p,
			New:    newChild,
		})
	}

	component.RC.Exec()

	return nil
}

// Update update component childes with error warning
func (component *Component) Update() {
	component.WarnError(component.UpdateWithError())
}

// ReCreate re create element
func (e *Element) ReCreate() {
	e.RC.Add(&RenderTask{
		Type:   RRecreate,
		New:    e,
		Parent: e.Parent,
	})
	go e.RC.Exec()
}

// UpdateChildes update element childes
func (e *Element) UpdateChildes() {
	e.OldChildes = e.Childes

	for i, childExt := range e.Childes {
		e.Childes[i] = e.RC.updateChild(e, childExt)
	}

	return
}

func (rc *RenderCore) updateChild(parent *Element, child interface{}) interface{} {
	if _, ok := child.(*Component); ok {
		childC := child.(*Component)
		childC.RC = rc
		child = childC.RenderElement()
	}

	childE, ok := child.(*Element)
	if !ok {
		return child
	}

	childE.RC = rc
	childE.Parent = parent

	if childE.Attrs != nil {
		childE.RAttrs = childE.Attrs()
	}

	childE.UpdateChildes()

	if childE.HTML != nil {
		childE.RHTML = childE.HTML()
	}

	return child
}

// UpdateElementChildes compare new and old trees
func (rc *RenderCore) UpdateElementChildes(_el interface{}, el *Element, new, old []interface{}, inReplaced bool) error {
	for i := 0; i < len(new) || i < len(old); i++ {
		var newEl interface{}
		if len(new) > i {
			newEl = new[i]
		}

		var oldEl interface{}
		if len(old) > i {
			oldEl = old[i]
		}

		err := rc.updateElement(_el, el, newEl, oldEl, i, inReplaced)
		if err != nil {
			return err
		}
	}

	return nil
}

// updateElement trying to update element
func (rc *RenderCore) updateElement(_parent interface{}, parent *Element, new, old interface{}, index int, inReplaced bool) error {
	// if element has created
	if old == nil {
		rc.Add(&RenderTask{
			Type:       RCreate,
			New:        new,
			Parent:     parent,
			NodeParent: _parent,

			InReplaced: inReplaced,
		})

		return nil
	}

	_childes := rc.BE.ChildNodes(_parent)

	var _el interface{}
	if len(_childes) > index {
		_el = _childes[index]
	} else {
		if IsElement(old) {
			_el = I2E(old).BEElement()
		} else {
			return nil
		}
	}

	// if element was deleted
	if new == nil {
		rc.Add(&RenderTask{
			Type:       RDelete,
			NodeParent: _parent,
			Parent:     parent,
			NodeOld:    _el,
			Old:        old,
			InReplaced: inReplaced,
		})

		return nil
	}

	// if element has Changed
	isChanged, canGoDeeper, err := Changed(new, old)
	if err != nil {
		return err
	}
	if isChanged {
		rc.Add(&RenderTask{
			Type:               RReplace,
			NodeParent:         _parent,
			NodeOld:            _el,
			Parent:             parent,
			New:                new,
			Old:                old,
			ReplaceCanGoDeeper: canGoDeeper,
			InReplaced:         inReplaced,
		})
	}
	if !canGoDeeper {
		return nil
	}

	newE := new.(*Element)
	oldE := old.(*Element)
	newE.UUID = oldE.UUID // little hack

	// if old and new is equals and they have html directives => they are two commons elements
	if oldE.HTML != nil && newE.HTML != nil {
		return nil
	}

	err = rc.UpdateElementChildes(_el, newE, newE.Childes, getOldChildes(newE, oldE), isChanged)
	if err != nil {
		return err
	}

	if isChanged && !inReplaced {
		rc.Add(&RenderTask{
			Type:       RReplaceHooks,
			NodeParent: _parent,
			NodeOld:    _el,
			Parent:     parent,
			New:        new,
			Old:        old,
		})
	}

	return nil
}

func getOldChildes(newE, oldE *Element) []interface{} {
	if newE.IsPointer {
		return newE.OldChildes
	} else {
		return oldE.Childes
	}
}
