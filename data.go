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

// SetData set data field and ForceUpdate component (after changes)
func (c *Component) SetData(query string, value interface{}) error {
	oldHtmlDirective := c.htmlDirective()

	oldValue := c.Data[query]
	err := c.SetDataFree(query, value)
	if err != nil {
		return err
	}

	if c.Watchers[query] != nil {
		err = c.Watchers[query](c, value, oldValue)
		if err != nil {
			return err
		}
	}

	return c.update(oldHtmlDirective)
}

// SetDataFree set data without ForceUpdate
func (c *Component) SetDataFree(query string, value interface{}) error {
	if value == nil {
		return fmt.Errorf("trying to set nil value to %s field", query)
	}

	c.Data[query] = value

	return nil
}

// DataDeleteFromArray remove element from data field
func (c *Component) DataDeleteFromArray(query string, index int) error {
	list, ok := c.GetData(query).([]interface{})
	if !ok {
		return errors.New("invalid data field type")
	}

	oldHtmlDirective := c.htmlDirective()

	err := c.SetDataFree(query, remove(list, index))
	if err != nil {
		return err
	}

	err = c.update(oldHtmlDirective)
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

	list = append(list, value)

	err := c.SetData(query, list)
	if err != nil {
		return err
	}

	return nil
}

// DataEditArray edit element in data field
func (c *Component) DataEditArray(query string, index int, value interface{}) error {
	list, ok := c.GetData(query).([]interface{})
	if !ok {
		return errors.New("invalid current list")
	}

	oldHtmlDirective := c.htmlDirective()

	list[index] = value

	err := c.update(oldHtmlDirective)
	if err != nil {
		return err
	}

	return nil
}

// remove remove item from element
func remove(a []interface{}, i int) []interface{} {
	copy(a[i:], a[i+1:]) // Shift a[i+1:] left one index
	a[len(a)-1] = ""     // Erase last element (write zero value)
	a = a[:len(a)-1]     // Truncate slice

	return a
}

// DataDeleteFromArray remove element from data field (works only with map[string]interface{} maps)
func (c *Component) DataDeleteFromMap(query string, key string) error {
	m, ok := c.GetData(query).(map[string]interface{})
	if !ok {
		return errors.New("invalid data field type")
	}

	oldHtmlDirective := c.htmlDirective()

	delete(m, key)

	err := c.update(oldHtmlDirective)
	if err != nil {
		return err
	}

	return nil
}

// DataEditMap edit element to data field (works only with map[string]interface{} maps)
func (c *Component) DataEditMap(query string, key string, value interface{}) error {
	m, ok := c.GetData(query).(map[string]interface{})
	if !ok {
		return errors.New("invalid data field type")
	}

	oldHtmlDirective := c.htmlDirective()

	m[key] = value

	err := c.update(oldHtmlDirective)
	if err != nil {
		return err
	}

	return nil
}
