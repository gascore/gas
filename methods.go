package gas

import (
	"errors"
	"fmt"
)

// Method runs a component method and updates component after
func (c *Component) Method(name string) error {
	method := c.Methods[name]
	if method == nil {
		return errors.New(fmt.Sprintf("invalid method name: %s", name))
	}

	oldTree := renderTree(c) // save rendered tree before

	err := method(c) // run method
	if err != nil {
		return err
	}

	newTree := renderTree(c) // save rendered tree after
	_c := c.GetElement()

	err = UpdateComponentChildes(_c, newTree, oldTree) // update component after changes
	if err != nil {
		return err
	}

	return nil
}

func (c *Component) Computed(name string) (interface{}, error) {
	computed := c.Computeds[name]
	if computed == nil {
		return nil, errors.New(fmt.Sprintf("invalid computed name: %s", name))
	}

	value, err := computed(c)
	if err != nil {
		return nil, err
	}

	return value, nil
}