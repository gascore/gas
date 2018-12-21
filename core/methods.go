package core

import (
	"errors"
	"fmt"
)

type PocketMethod func(...interface{}) error

type PocketComputed func(...interface{})(interface{}, error)

// Method runs a component method and updates component after
func (c *Component) Method(name string, values ...interface{}) error {
	method, err := c.GetPocketMethod(name)
	if err != nil {
		return errors.New(fmt.Sprintf("invalid method name: %s", name))
	}

	err = method(values...) // run method
	if err != nil {
		return err
	}

	return nil
}

// GetPocketMethod return function returns executing method with binding component
func (c *Component) GetPocketMethod(name string) (PocketMethod, error)  {
	method := c.Methods[name]
	if method == nil {
		return nil, errors.New(fmt.Sprintf("invalid method name: %s", name))
	}

	bindingMethod := func(values ...interface{}) error {
		return method(c, values...)
	}

	return bindingMethod, nil
}

// Computed runs a component computed and returns values from it
func (c *Component) Computed(name string, values ...interface{}) (interface{}, error) {
	computed, err := c.GetPocketComputed(name)
	if err != nil {
		return nil, err
	}

	value, err := computed(values...)
	if err != nil {
		return nil, err
	}

	return value, nil
}

// GetPocketComputed return function returns executing computed with binding component
func (c *Component) GetPocketComputed(name string) (PocketComputed, error)  {
	computed := c.Computeds[name]
	if computed == nil {
		return nil, errors.New(fmt.Sprintf("invalid computed name: %s", name))
	}

	bindingComputed := func(values ...interface{}) (interface{}, error) {
		return computed(c, values...)
	}

	return bindingComputed, nil
}