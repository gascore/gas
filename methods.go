package gas

import (
	"fmt"
	"github.com/pkg/errors"
)

type PocketMethod func(...interface{}) (interface{}, error)

type PocketComputed func(...interface{}) (interface{}, error)

// Method runs a component method and updates component after
func (c *Component) Method(name string, values ...interface{}) interface{} {
	out, err := c.MethodSafely(name, values...)
	if err != nil {
		c.ConsoleError(err.Error())
		return nil
	}

	return out
}

func (c *Component) MethodSafely(name string, values ...interface{}) (interface{}, error) {
	method := c.PocketMethod(name)
	if method == nil {
		return nil, fmt.Errorf("invalid method name: %s", name)
	}

	return method(values...)
}

// PocketMethod return function returns executing method with binding component
func (c *Component) PocketMethod(name string) PocketMethod {
	method := c.Methods[name]
	if method == nil {
		c.WarnError(fmt.Errorf("invalid method name: %s", name))
		return nil
	}

	return func(values ...interface{}) (interface{}, error) {
		return method(c, values...)
	}
}

// TODO: Add caching for computeds

// Computed runs a component computed and returns values from it
func (c *Component) Computed(name string, values ...interface{}) interface{} {
	out, err := c.ComputedSafely(name, values...)
	if err != nil {
		c.ConsoleError(err.Error())
		return nil
	}

	return out
}

func (c *Component) ComputedSafely(name string, values ...interface{}) (interface{}, error) {
	computed := c.PocketComputed(name)
	if computed == nil {
		return nil, errors.New("invalid computed: Computeds[name] is nil")
	}

	return computed(values...)
}

// PocketComputed return function returns executing computed with binding component
func (c *Component) PocketComputed(name string) PocketComputed {
	computed := c.Computeds[name]
	if computed == nil {
		c.WarnError(fmt.Errorf("invalid computed name: %s", name))
		return nil
	}

	return func(values ...interface{}) (interface{}, error) {
		val, err := computed(c, values...)
		if err != nil {
			return val, err
		}

		return val, nil
	}
}
