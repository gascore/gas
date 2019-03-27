package gas

import (
	"fmt"
)

// PocketMethod executeble component Method passing to non-component-context places (to component childes, etc)
type PocketMethod func(...interface{}) (interface{}, error)

// Method runs a component method and updates component after
func (c *Component) Method(name string, values ...interface{}) interface{} {
	out, err := c.MethodSafely(name, values...)
	if err != nil {
		c.ConsoleError(err.Error())
		return nil
	}

	return out
}

// MethodSafely call component Method safely with returning error
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
