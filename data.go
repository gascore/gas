package gas

import (
	"fmt"
	"github.com/Sinicablyat/dom"
	"github.com/pkg/errors"
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
		dom.ConsoleError(fmt.Sprintf("trying to set nil value to %s field", query))
	}

	c.Data[query] = value

	return nil
}


// getOldValues return `old` values for component update
func (c *Component) getOldValues() ([]interface{}, string) {
	oldTree := renderTree(c)

	var oldHtmlDirective string
	if c.Directives.HTML.Render != nil {
		oldHtmlDirective = c.Directives.HTML.Render(c)
	}

	return oldTree, oldHtmlDirective
}

// DoWithUpdate runs your event and trying to Update component after it
func (c *Component) DoWithUpdate(event func()error) error {
	oldTree, oldHtmlDirective := c.getOldValues()

	err := event() // your event
	if err != nil {
		return err
	}

	return c.Update(oldTree, oldHtmlDirective)
}

// Update update component
func (c *Component) Update(oldTree []interface{}, oldHtmlDirective string) error {
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

	err := UpdateComponentChildes(_c, newTree, oldTree)
	if err != nil {
		return err
	}

	return nil
}


// DataDeleteFromArray Remove element from data field
func (c *Component) DataDeleteFromArray(query string, index int) error {
	list, ok := c.GetData(query).([]interface{})
	if !ok {
		return errors.New("invalid data field type")
	}

	oldTree, oldHtmlDirective := c.getOldValues()

	err := c.SetDataFree(query, Remove(list, index))
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

	oldTree, oldHtmlDirective := c.getOldValues()


	list[index] = value


	err := c.Update(oldTree, oldHtmlDirective)
	if err != nil {
		return err
	}

	return nil
}


// Remove remove item from element
func Remove(a []interface{}, i int) []interface{} {
	copy(a[i:], a[i+1:]) // Shift a[i+1:] left one index
	a[len(a)-1] = ""     // Erase last element (write zero value)
	a = a[:len(a)-1]     // Truncate slice

	return a
}