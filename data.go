package gas

import (
	"errors"
	"fmt"
)

// GetData return data field by query string
func (c *Component) GetData(query string) interface{} {
	// There will be callbacks, events, e.t.c.
	data := c.Data[query]

	return data
}

// SetData set data field and Update component (after changes)
func (c *Component) SetData(query string, value interface{}) error {
	err := c.DoWithUpdate(func() error {
		oldValue := c.Data[query]
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

// SetDataFree set data without Update
func (c *Component) SetDataFree(query string, value interface{}) error {
	if value == nil {
		return fmt.Errorf("trying to set nil value to %s field", query)
	}

	c.Data[query] = value

	return nil
}


// getOldState return *old* values for component update
func (c *Component) getOldState() ([]interface{}, string) {
	return c.RChildes, c.renderHtmlDirective()
}

// getState return values for component update
func (c *Component) getState() ([]interface{}, string) {
	return c.be.RenderTree(c), c.renderHtmlDirective()
}

func (c *Component) renderHtmlDirective() string {
	var htmlDirective string
	if c.Directives.HTML.Render != nil {
		htmlDirective = c.Directives.HTML.Render(c)
	}

	return htmlDirective
}

// DoWithUpdate runs your event and trying to Update component after it
func (c *Component) DoWithUpdate(event func()error) error {
	oldTree, oldHtmlDirective := c.getState()

	err := event() // your event
	if err != nil {
		return err
	}

	return c.Update(oldTree, oldHtmlDirective)
}

// Update update component
func (c *Component) Update(oldTree []interface{}, oldHtmlDirective string) error {
	newTree, newHtmlDirective := c.getState()

	if oldHtmlDirective != newHtmlDirective {
		err := c.be.ReCreate(c)
		if err != nil {
			return err
		}

		return nil
	}

	err := c.be.UpdateComponentChildes(c, newTree, oldTree)
	if err != nil {
		return err
	}

	return nil
}

// ReCreate recreate your component
func (c *Component) ReCreate() error {
	return c.be.ReCreate(c)
}

// ReCreate recreate your component
func (c *Component) Reload() error {
	return c.be.ReloadComponent(c)
}

// ForceUpdate force update your component
func (c *Component) ForceUpdate() error {
	childes := c.Childes(c)
	err := c.be.UpdateComponentChildes(c, childes, c.RChildes)
	if err != nil {
		return err
	}

	return nil
}


// DataDeleteFromArray remove element from data field
func (c *Component) DataDeleteFromArray(query string, index int) error {
	list, ok := c.GetData(query).([]interface{})
	if !ok {
		return errors.New("invalid data field type")
	}

	oldTree, oldHtmlDirective := c.getState()

	err := c.SetDataFree(query, remove(list, index))
	if err != nil {
		return err
	}

	err = c.Update(oldTree, oldHtmlDirective)
	if err != nil {
		return err
	}

	return nil
}

// DataAddToArray add element to data field
func (c *Component) DataAddToArray(query string, value interface{}) error {
	list, ok := c.GetData(query).([]interface{})
	if !ok {
		return errors.New("invalid data field type")
	}

	list  = append(list, value)

	err := c.SetData(query, list)
	if err != nil {
		return err
	}

	return nil
}

// DataEditArray edit element in data field
func (c *Component) DataEditArray(query string, index int, value interface{}) error {
	list, ok := c.GetData("current").([]interface{})
	if !ok {
		return errors.New("invalid current list")
	}

	oldTree, oldHtmlDirective := c.getState()

	list[index] = value

	err := c.Update(oldTree, oldHtmlDirective)
	if err != nil {
		return err
	}

	return nil
}

// TODO: Add methods for work with map[string]interface{}

// remove remove item from element
func remove(a []interface{}, i int) []interface{} {
	copy(a[i:], a[i+1:]) // Shift a[i+1:] left one index
	a[len(a)-1] = ""     // Erase last element (write zero value)
	a = a[:len(a)-1]     // Truncate slice

	return a
}