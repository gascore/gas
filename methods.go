package gas

import (
	"errors"
	"fmt"
)

type PocketMethod func(...interface{}) error

type PocketComputed func(...interface{})interface{}

// Method runs a component method and updates component after
func (c *Component) Method(name string, values ...interface{}) error {
	method := c.GetPocketMethod(name)

	err := method(values...) // run method
	if err != nil {
		return err
	}

	return nil
}

// GetPocketMethod return function returns executing method with binding component
func (c *Component) GetPocketMethod(name string) PocketMethod  {
	method := c.Methods[name]
	if method == nil {
		WarnError(errors.New(fmt.Sprintf("invalid method name: %s", name)))
		return nil
	}

	bindingMethod := func(values ...interface{}) error {
		return method(c, values...)
	}

	return bindingMethod
}

// Computed runs a component computed and returns values from it
func (c *Component) Computed(name string, values ...interface{}) interface{} {
	computed := c.GetPocketComputed(name)

	value := computed(values...)

	return value
}

// GetPocketComputed return function returns executing computed with binding component
func (c *Component) GetPocketComputed(name string) PocketComputed  {
	computed := c.Computeds[name]
	if computed == nil {
		WarnError(errors.New(fmt.Sprintf("invalid computed name: %s", name)))
		return nil
	}

	bindingComputed := func(values ...interface{}) interface{} {
		val, err := computed(c, values...)
		if err != nil {
			WarnError(err)
			return nil
		}

		return val
	}

	return bindingComputed
}