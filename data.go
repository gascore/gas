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
	oldValue := c.Data[query]
	err := c.eventInUpdater(func() error {
		err := c.SetDataFree(query, value)
		if err != nil {
			return err
		}

		if c.Watchers[query] == nil {
			return nil
		}

		err = c.Watchers[query](c, value, oldValue)
		if err != nil {
			return err
		}

		return nil
	})
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

// eventInUpdater runs your event and trying to update component after it
func (c *Component) eventInUpdater(event func()error) error {
	oldTree := renderTree(c)

	var oldHtmlDirective string
	if c.Directives.HTML.Render != nil {
		oldHtmlDirective = c.Directives.HTML.Render(c)
	}

	err := event() // your event
	if err != nil {
		return err
	}

	newTree := renderTree(c)

	var newHtmlDirective string
	if c.Directives.HTML.Render != nil {
		newHtmlDirective = c.Directives.HTML.Render(c)
	}

	_c := c.GetElement()

	if oldHtmlDirective != newHtmlDirective {
		_updatedC, err := CreateComponent(c)
		if err != nil {
			return err
		}

		c.ParentC.GetElement().ReplaceChild(_updatedC, _c)
		return nil
	}

	err = UpdateComponentChildes(_c, newTree, oldTree)
	if err != nil {
		return err
	}

	return nil
}