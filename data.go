package gas

import (
	"errors"
)

var (
	// ErrComponentDataIsNil if component.Data == nil
	ErrComponentDataIsNil = errors.New("component Data is nil")

	// ErrNilField if field == nil
	ErrNilField = errors.New("trying to set value to nil field")

	// ErrInvalidDataFieldType if fieldValue.(type) != newFieldValue.(type)
	ErrInvalidDataFieldType = errors.New("invalid Data field type")
)

// Get return Data field by query string
func (c *Component) Get(query string) interface{} {
	if c.Data == nil {
		c.ConsoleError(ErrComponentDataIsNil.Error())
		return nil
	}

	if _, ok := c.Data[query]; !ok {
		c.ConsoleError("field is undefined")
		return nil
	}

	data := c.Data[query]

	return data
}

// Set set many values for many Data fields and ForceUpdate component
func (c *Component) Set(data map[string]interface{}) {
	if data == nil {
		return
	}

	c.RC.Add(singleNode(&RenderNode{
		Type:     DataType,
		Priority: EventPriority,

		New:  c,
		Data: data,
	}))
}

// SetValue set Data field and ForceUpdate component
func (c *Component) SetValue(query string, value interface{}) {
	c.Set(map[string]interface{}{query: value})
}

// Set set many values for many Data fields and ForceUpdate component
func (c *Component) realSet(node *RenderNode) error {
	oldHTMLDirective := c.htmlDirective()

	if node.Data == nil {
		return errors.New("invalid Data for Set")
	}

	err := c.SetImm(node.Data)
	if err != nil {
		return err
	}

	return c.update(oldHTMLDirective)
}

// SetValueImm set data immediately + without ForceUpdate
func (c *Component) SetValueImm(query string, value interface{}) error {
	if c.Data == nil {
		c.Data = make(map[string]interface{})
		return ErrComponentDataIsNil
	}

	if _, ok := c.Data[query]; !ok {
		return ErrNilField
	}

	oldValue := c.Data[query]
	c.Data[query] = value

	if c.Watchers[query] != nil {
		err := c.Watchers[query](c, value, oldValue)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetImm SetValueImm by map
func (c *Component) SetImm(data map[string]interface{}) error {
	for key, value := range data {
		err := c.SetValueImm(key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// DataDeleteFromArray remove element from Data field
func (c *Component) DataDeleteFromArray(query string, index int) error {
	list, ok := c.Get(query).([]interface{})
	if !ok {
		return ErrInvalidDataFieldType
	}

	oldHTMLDirective := c.htmlDirective()

	err := c.SetValueImm(query, remove(list, index))
	if err != nil {
		return err
	}

	err = c.update(oldHTMLDirective)
	if err != nil {
		return err
	}

	return nil
}

// DataAddToArray add element to Data field
func (c *Component) DataAddToArray(query string, value interface{}) error {
	list, ok := c.Get(query).([]interface{})
	if !ok {
		return ErrInvalidDataFieldType
	}

	list = append(list, value)

	c.SetValue(query, list)

	return nil
}

// DataEditArray edit element in Data field
func (c *Component) DataEditArray(query string, index int, value interface{}) error {
	list, ok := c.Get(query).([]interface{})
	if !ok {
		return ErrInvalidDataFieldType
	}

	oldHTMLDirective := c.htmlDirective()

	list[index] = value

	err := c.update(oldHTMLDirective)
	if err != nil {
		return err
	}

	return nil
}

// DataDeleteFromMap remove element from Data field (works only with map[string]interface{} maps)
func (c *Component) DataDeleteFromMap(query, key string) error {
	m, ok := c.Get(query).(map[string]interface{})
	if !ok {
		return ErrInvalidDataFieldType
	}

	oldHTMLDirective := c.htmlDirective()

	delete(m, key)

	err := c.update(oldHTMLDirective)
	if err != nil {
		return err
	}

	return nil
}

// DataEditMap edit element to Data field (works only with map[string]interface{} maps)
func (c *Component) DataEditMap(query, key string, value interface{}) error {
	m, ok := c.Get(query).(map[string]interface{})
	if !ok {
		return ErrInvalidDataFieldType
	}

	oldHTMLDirective := c.htmlDirective()

	m[key] = value

	err := c.update(oldHTMLDirective)
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
