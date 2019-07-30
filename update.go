package gas

import (
	"errors"
)

// htmlDirective return compiled element HTMLDirective
func (e *Element) htmlDirective() string {
	var htmlDirective string
	if e.HTML.Render != nil {
		htmlDirective = e.HTML.Render()
	}

	return htmlDirective
}

func (e *Element) Update() error {
	_el := e.BEElement()
	if _el == nil {
		return errors.New("invalid '_el' in update function")
	}

	e.UpdateChildes()
	err := e.RC.UpdateElementChildes(_el, e, e.Childes, e.OldChildes)
	if err != nil {
		return err
	}

	return e.RC.Exec()
}

// UpdateHTMLDirective trying rerender element html directive
func (e *Element) UpdateHTMLDirective() {
	if e.HTML.Render != nil {
		e.HTML.Rendered = e.HTML.Render()
	}
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
	e.RC.Add(&RenderNode{
		Type: RecreateType,
		New:  e,
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
		childE.UpdateHTMLDirective()
	}

	return
}

// UpdateElementChildes compare new and old trees
func (rc *RenderCore) UpdateElementChildes(_el interface{}, el *Element, new, old []interface{}) error {
	for i := 0; i < len(new) || i < len(old); i++ {
		var newEl interface{}
		if len(new) > i {
			newEl = new[i]
		}

		var oldEl interface{}
		if len(old) > i {
			oldEl = old[i]
		}

		err := rc.updateElement(_el, el, newEl, oldEl, i)
		if err != nil {
			return err
		}
	}

	return nil
}

// updateElement trying to update element
func (rc *RenderCore) updateElement(_parent interface{}, parent *Element, new, old interface{}, index int) error {
	// if element has created
	if old == nil {
		rc.Add(&RenderNode{
			Type:       CreateType,
			New:        new,
			Parent: 	parent,
			NodeParent: _parent,
		})

		return nil
	}

	_childes := rc.BE.ChildNodes(_parent)
	if _childes == nil {
		_childes = []interface{}{}
	}

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
		rc.Add(&RenderNode{
			Type:       DeleteType,
			NodeParent: _parent,
			Parent: 	parent,
			NodeOld:    _el,
			Old:        old,
		})

		return nil
	}

	// if element has Changed
	isChanged, err := Changed(new, old)
	if err != nil {
		return err
	}

	if isChanged {
		rc.Add(&RenderNode{
			Type:       ReplaceType,
			NodeParent: _parent,
			NodeOld:    _el,
			Parent: 	parent,
			New:        new,
			Old:        old,
		})

		return nil
	}

	newE, newIsElement := new.(*Element)
	if !newIsElement {
		return nil
	}

	oldE := old.(*Element)

	if newE.UUID != oldE.UUID { // update *element context* in new
		newE.UUID = oldE.UUID
	}

	// if old and new is equals and they have html directives => they are two commons elements
	if oldE.HTML.Render != nil && newE.HTML.Render != nil {
		return nil
	}

	if newE.IsPointer {
		err = rc.UpdateElementChildes(_el, newE, newE.Childes, newE.OldChildes)
	} else {
		err = rc.UpdateElementChildes(_el, oldE, newE.Childes, oldE.Childes)
	}
	if err != nil {
		return err
	}

	return nil
}
