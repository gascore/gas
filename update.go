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
	newTree := c.be.RenderTree(c)

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

	return nil
}

// ForceUpdate force update component
func (c *Component) ForceUpdate() error {
	return c.update("")
}