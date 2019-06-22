package gas

import (
	"errors"
	"sync"
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
	newTree := e.RenderTree()

	if e.HTML.Rendered != e.htmlDirective() {
		e.ReCreate()
		return nil
	}

	_el := e.BEElement()
	if _el == nil {
		return errors.New("invalid '_el' in update function")
	}

	renderNodes, err := e.RC.UpdateElementChildes(_el, newTree, e.RChildes)
	if err != nil {
		return err
	}

	e.RChildes = newTree
	e.UpdateHTMLDirective()

	e.RC.Add(renderNodes)

	return nil
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
	e.RC.Add(singleNode(&RenderNode{
		Type: RecreateType,
		New:  e,
	}))
}

// RenderTree return full rendered childes tree of element
func (e *Element) RenderTree() []interface{} {
	var childes []interface{}
	for _, child := range UnSpliceBody(e.Childes()) {
		childE, ok := child.(*Element)
		if !ok {
			childes = append(childes, child)
			continue
		}

		childE.RC = e.RC
		childE.Parent = e

		if childE.Binds != nil {
			if childE.RenderedBinds == nil {
				childE.RenderedBinds = make(map[string]string)
			}

			// render binded attributes
			for bindKey, bindValue := range childE.Binds {
				childE.RenderedBinds[bindKey] = bindValue()
			}
		}

		childE.RChildes = childE.RenderTree()
		childE.UpdateHTMLDirective()

		child = childE
		childes = append(childes, child)
	}

	return childes
}

// UpdateElementChildes compare new and old trees
func (rc *RenderCore) UpdateElementChildes(_el interface{}, newTree, oldTree []interface{}) ([]*RenderNode, error) {
	var nodes []*RenderNode

	var m sync.Mutex
	var errG error

	for i := 0; i < len(newTree) || i < len(oldTree); i++ {
		var newEl interface{}
		if len(newTree) > i {
			newEl = newTree[i]
		}

		var oldEl interface{}
		if len(oldTree) > i {
			oldEl = oldTree[i]
		}

		renderNodes, err := rc.updateElement(_el, newEl, oldEl, i)
		if err != nil {
			m.Lock()
			errG = err
			m.Unlock()
		}

		if renderNodes != nil {
			m.Lock()
			nodes = append(nodes, renderNodes...)
			m.Unlock()
		}
	}

	return nodes, errG
}

// updateElement trying to update element
func (rc *RenderCore) updateElement(_parent interface{}, new interface{}, old interface{}, index int) ([]*RenderNode, error) {
	var nodes []*RenderNode

	// if element has created
	if old == nil {
		nodes = append(nodes, &RenderNode{
			Type:       CreateType,
			New:        new,
			NodeParent: _parent,
		})

		return nodes, nil
	}

	_childes := rc.BE.ChildNodes(_parent)
	if _childes == nil {
		_childes = []interface{}{}
	}

	var _el interface{}
	if len(_childes) > index { // element was hided if childes length <= index
		_el = _childes[index]
	} else {
		return nodes, nil
	}

	newIsElement := IsElement(new)
	newE := &Element{}
	if newIsElement {
		newE = I2E(new)
	}

	// if element has deleted
	if new == nil {
		nodes = append(nodes, &RenderNode{
			Type:       DeleteType,
			NodeParent: _parent,
			NodeOld:    _el,
			Old:        old,
		})

		return nodes, nil
	}

	// if element has Changed
	isChanged, err := Changed(new, old)
	if err != nil {
		return nil, err
	}
	if isChanged {
		nodes = append(nodes, &RenderNode{
			Type:       ReplaceType,
			NodeParent: _parent,
			NodeOld:    _el,
			New:        new,
			Old:        old,
		})

		return nodes, nil
	}

	if !newIsElement {
		return nodes, nil
	}

	oldE := I2E(old)

	if newE.UUID != oldE.UUID { // update *element context* in new
		newE.UUID = oldE.UUID
	}

	// if old and new is equals and they have html directives => they are two commons elements
	if oldE.HTML.Render != nil && newE.HTML.Render != nil {
		return nodes, nil
	}

	if newE.Component != nil {
		newE.RChildes = newE.RenderTree()
	}

	renderNodes, err := rc.UpdateElementChildes(_el, newE.RChildes, oldE.RChildes)
	if err != nil {
		return nil, err
	}

	nodes = append(nodes, renderNodes...)

	return nodes, nil
}


// UpdateWatchersValues update all elements WatcherValue where element.Watcher = watcher 
func (component *Component) UpdateWatchersValues(watcher string, newVal string) {
	component.Element.updateWatchersValues(watcher, newVal)
}

func (e *Element) updateWatchersValues(watcher string, newVal string) {
	for _, child := range e.RChildes {
		childE, ok := child.(*E)
		if !ok {
			continue
		}

		if childE.Watcher == watcher {
			e.RC.BE.EditWatcherValue(childE.BEElement(), newVal)
			continue
		}

		childE.updateWatchersValues(watcher, newVal)
	}
}
