package gas

import (
	"fmt"
	"github.com/Sinicablyat/dom"
)

// GetData return data field by query string
func (c *Component) GetData(query string) interface{} {
	// There will be callbacks, events, e.t.c.
	data := c.Data[query]

	return data
}

// SetData set data field and update component (after changes)
func (c *Component) SetData(query string, value interface{}) error {
	oldTree := renderTree(c)

	_ = c.SetDataFree(query, value)

	newTree := renderTree(c)
	_c := c.GetElement()

	err := UpdateComponentChildes(_c, newTree, oldTree)
	if err != nil {
		return err
	}

	return nil
}

// SetDataFree set data without update
func (c *Component) SetDataFree(query string, value interface{}) error {
	if value == nil {
		dom.ConsoleError(fmt.Sprintf("trying to set nil value to %s field", query))
	}

	c.Data[query] = value

	return nil
}