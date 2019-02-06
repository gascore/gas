package gas

// htmlDirective return compiled component HTMLDirective
func (c *Component) htmlDirective() string {
	var htmlDirective string
	if c.Directives.HTML.Render != nil {
		htmlDirective = c.Directives.HTML.Render(c)
	}

	return htmlDirective
}

func (c *Component) update(oldHtmlDirective string) error {
	newTree := RenderTree(c)

	if oldHtmlDirective != c.htmlDirective() {
		err := c.be.ReCreate(c)
		if err != nil {
			return err
		}

		return nil
	}

	err := c.be.UpdateComponentChildes(c, newTree, c.RChildes)
	if err != nil {
		return err
	}

	c.RChildes = newTree
	c.UpdateHtmlDirective()

	return nil
}

// UpdateHtmlDirective trying rerender component html directive
func (c *Component) UpdateHtmlDirective() {
	if c.Directives.HTML.Render != nil {
		c.Directives.HTML.Rendered = c.Directives.HTML.Render(c)
	}
}

// ForceUpdate force update component
func (c *Component) ForceUpdate() error {
	return c.update(c.Directives.HTML.Rendered)
}

func (c *Component) ReCreate() error {
	return c.be.ReCreate(c)
}

// RenderTree return full rendered childes tree of component
func RenderTree(c *Component) []interface{} {
	var childes []interface{}
	for _, el := range c.Childes(c) {
		if IsComponent(el) {
			elC := I2C(el)

			if elC.Binds != nil {
				if elC.RenderedBinds == nil {
					elC.RenderedBinds = map[string]string{}
				}

				for bindKey, bindValue := range elC.Binds { // render binds
					elC.RenderedBinds[bindKey] = bindValue()
				}
			}

			elC.RChildes = RenderTree(elC)
			elC.UpdateHtmlDirective()

			el = elC
		}

		childes = append(childes, el)
	}
	return childes
}
