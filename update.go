package gas

import (
	"errors"
)

// htmlDirective return compiled component HTMLDirective
func (c *Component) htmlDirective() string {
	var htmlDirective string
	if c.HTML.Render != nil {
		htmlDirective = c.HTML.Render(c)
	}

	return htmlDirective
}

func (c *Component) update(oldHTMLDirective string) error {
	newTree := RenderTree(c)

	if oldHTMLDirective != c.htmlDirective() {
		c.ReCreate()
		return nil
	}

	_el := c.Element()
	if _el == nil {
		return errors.New("invalid '_el' value in DeepUpdateComponentChildes")
	}

	renderNodes, err := UpdateComponentChildes(c, _el, newTree, c.RChildes)
	if err != nil {
		return err
	}

	c.RChildes = newTree
	c.UpdateHTMLDirective()

	c.RC.Add(renderNodes)

	return nil
}

// UpdateHTMLDirective trying rerender component html directive
func (c *Component) UpdateHTMLDirective() {
	if c.HTML.Render != nil {
		c.HTML.Rendered = c.HTML.Render(c)
	}
}

// ForceUpdate force update component
func (c *Component) ForceUpdate() error {
	return c.update(c.HTML.Rendered)
}

// ReCreate re create component
func (c *Component) ReCreate() {
	c.RC.Add(singleNode(&RenderNode{
		Type: RecreateType,
		Priority: RenderPriority,
		New: c,
	}))
}

// RenderTree return full rendered childes tree of component
func RenderTree(c *Component) []interface{} {
	var childes []interface{}
	for _, el := range c.Childes(c) {
		if IsComponent(el) {
			elC := I2C(el)

			if elC.Binds != nil {
				if elC.RenderedBinds == nil {
					elC.RenderedBinds = make(map[string]string)
				}

				for bindKey, bindValue := range elC.Binds { // render binds
					elC.RenderedBinds[bindKey] = bindValue()
				}
			}

			elC.RChildes = RenderTree(elC)
			elC.UpdateHTMLDirective()

			el = elC
		}

		childes = append(childes, el)
	}

	return childes
}

func UpdateComponentChildes(c *Component, _el interface{}, newTree, oldTree []interface{}) ([]*RenderNode, error) {
	var nodes []*RenderNode

	for i := 0; i < len(newTree) || i < len(oldTree); i++ {
		var elFromNew interface{}
		if len(newTree) > i {
			elFromNew = newTree[i]
		}

		var elFromOld interface{}
		if len(oldTree) > i {
			elFromOld = oldTree[i]
		}

		renderNodes, err := c.RC.updateComponent(_el, elFromNew, elFromOld, i)
		if err != nil {
			return nil, err
		}

		if renderNodes != nil {
			nodes = append(nodes, renderNodes...)
		}
	}

	return nodes, nil
}

// updateComponent trying to update component
func (rc *RenderCore) updateComponent(_parent interface{}, new interface{}, old interface{}, index int) ([]*RenderNode, error) {
	var nodes []*RenderNode
	// if component has created
	if old == nil {
		nodes = append(nodes, &RenderNode{
			Type: CreateType,
			Priority: RenderPriority,
			New: new,
			NodeParent: _parent,
		})

		return nodes, nil
	}

	_childes := rc.BE.ChildNodes(_parent)
	if _childes == nil {
		_childes = []interface{}{}
	}

	var _el interface{}
	if len(_childes) > index { // component was hided if childes length <= index
		_el = _childes[index]
	} else {
		return nodes, nil
	}

	newIsComponent := IsComponent(new)
	newC := &Component{}
	if newIsComponent {
		newC = I2C(new)
	}

	// if component has deleted
	if new == nil {
		nodes = append(nodes, &RenderNode{
			Type:DeleteType,
			Priority: RenderPriority,
			NodeParent:_parent,
			NodeOld: _el,
			Old: old,
		})

		return nodes, nil
	}

	// if component has Changed
	isChanged, err := Changed(new, old)
	if err != nil {
		return nil, err
	}
	if isChanged {
		nodes = append(nodes, &RenderNode{
			Type:ReplaceType,
			Priority: RenderPriority,
			NodeParent:_parent,
			NodeOld:_el,
			New: new,
			Old: old,
		})

		return nodes, nil
	}

	if !newIsComponent {
		return nodes, nil
	}

	oldC := I2C(old)

	if newC.UUID != oldC.UUID { // update *element context* in new
		newC.UUID = oldC.UUID
	}

	// if old and new is equals and they have html directives => they are two commons components
	if IsComponent(old) && oldC.HTML.Render != nil && newC.HTML.Render != nil {
		return nodes, nil
	}

	nodes = append(nodes, &RenderNode{
		Type: SyncType,
		Priority: RenderPriority,
		New: newC,
		NodeNew: _el,
	})

	renderNodes, err := UpdateComponentChildes(newC, _el, newC.RChildes, oldC.RChildes)
	if err != nil {
		return nil, err
	}

	nodes = append(nodes, renderNodes...)

	return nodes, nil
}
