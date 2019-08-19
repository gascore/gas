package gas

import (
	"errors"
)

func (e *Element) Update() error {
	_el := e.BEElement()
	if _el == nil {
		return errors.New("invalid '_el' in update function")
	}

	e.UpdateChildes()
	err := e.RC.UpdateElementChildes(_el, e, e.Childes, e.OldChildes, false)
	if err != nil {
		return err
	}

	e.RC.Exec()

	return nil
}

// UpdateWithError update component childes
func (component *Component) UpdateWithError() error {
	return component.Element.Update()
}

// Update update component childes with error warning
func (component *Component) Update() {
	component.WarnError(component.UpdateWithError())
}

// ReCreate re create element
func (e *Element) ReCreate() {
	e.RC.Add(&RenderTask{
		Type:   RecreateType,
		New:    e,
		Parent: e.Parent,
	})
	go e.RC.Exec()
}

// UpdateChildes update element childes
func (e *Element) UpdateChildes() {
	e.OldChildes = e.Childes
	e.Childes = []interface{}{}

	for _, child := range UnSpliceBody(e.getChildes()) {
		e.Childes = append(e.Childes, child)

		childE, ok := child.(*Element)
		if !ok {
			continue
		}

		childE.RC = e.RC
		childE.Parent = e

		if childE.Attrs != nil {
			childE.RAttrs = childE.Attrs()
		}

		childE.UpdateChildes()

		if childE.HTML != nil {
			childE.RHTML = childE.HTML()
		}
	}

	return
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
			Type:        CreateType,
			New:         new,
			Parent:      parent,
			NodeParent:  _parent,
			IgnoreHooks: inReplaced,
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
			Type:        DeleteType,
			NodeParent:  _parent,
			Parent:      parent,
			NodeOld:     _el,
			Old:         old,
			IgnoreHooks: inReplaced,
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
			Type:               ReplaceType,
			NodeParent:         _parent,
			NodeOld:            _el,
			Parent:             parent,
			New:                new,
			Old:                old,
			ReplaceCanGoDeeper: canGoDeeper,
			IgnoreHooks:        inReplaced,
		})
	}
	if !canGoDeeper {
		return nil
	}

	newE := new.(*Element)
	oldE := old.(*Element)
	if newE.UUID != oldE.UUID {
		newE.UUID = oldE.UUID // little hack
	}

	// if old and new is equals and they have html directives => they are two commons elements
	if oldE.HTML != nil && newE.HTML != nil {
		return nil
	}

	var oldChildes []interface{}
	if newE.IsPointer {
		oldChildes = newE.OldChildes
	} else {
		oldChildes = oldE.Childes
	}

	err = rc.UpdateElementChildes(_el, newE, newE.Childes, oldChildes, isChanged)
	if err != nil {
		return err
	}

	if isChanged && !inReplaced {
		rc.Add(&RenderTask{
			Type:       ReplaceHooks,
			NodeParent: _parent,
			NodeOld:    _el,
			Parent:     parent,
			New:        new,
			Old:        old,
		})
	}

	return nil
}
