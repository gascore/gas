package gas

import (
	"fmt"
	"github.com/Sinicablyat/dom"
)

// GetData return data field by query string
func (c *Component) GetData(query string) interface{} {
	// There will be callbacks, events, e.t.c.
	data := c.Data[query]
	if data == nil {
		dom.ConsoleError(fmt.Sprintf(`"%s"trying to accept nil data`, c.Tag))
	}

	return data
}

// SetData set data field and update component (after changes)
func (c *Component) SetData(query string, value interface{}) error {
	oldTree := renderTree(c)


	c.Data[query] = value


	newTree := renderTree(c)
	_c := c.GetElement()

	err := UpdateComponentChildes(_c, newTree, oldTree)
	if err != nil {
		return err
	}

	return nil
}