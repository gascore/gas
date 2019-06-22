package gas

import "github.com/frankenbeanies/uuid4"

// Component logic node
type Component struct {
	Root  interface {
		Render() []interface{}
	}

	Element *Element

	Hooks Hooks

	RefsAllowed bool
	Refs        map[string]*Element

	Watchers map[string]Watcher

	RC *RenderCore
}

// Watcher function called when input event triggering
type Watcher func(val interface{}, e Object) (string, error)

// Init initialize component: create element and other stuff
func (c *Component) Init() *Element {
	el := &Element{
		Tag: "div",
		UUID: uuid4.New().String(),
		Component: c,
		RC: c.RC,
	}
	
	el.Childes = func() []interface{} {
		return c.Root.Render()
		/*var compiled []interface{}
		for _, child := range UnSpliceBody(c.Root.Render()) {
			childE, ok := child.(*Element)
			if !ok {
				compiled = append(compiled, child)
				continue
			}

			childE.RC = c.RC
			childE.Parent = el

			compiled = append(compiled, child)
		}

		return compiled */
	}

	c.Element = el

	return el
}

// UnSpliceBody extract values fromm array to component childes
func UnSpliceBody(body []interface{}) []interface{} {
	var arr []interface{}
	for _, el := range body {
		switch el.(type) {
		case []interface{}:
			for _, c := range el.([]interface{}) {
				arr = append(arr, c)
			}
			continue
		case []*E:
			for _, e := range el.([]*E) {
				arr = append(arr, e)
			}
			continue
		case []*C:
			for _, e := range el.([]*C) {
				arr = append(arr, e)
			}
			continue
		case nil:
			continue
		default:
			arr = append(arr, el)
		}
	}

	return arr
}

// I2C - convert interface{} to *Component
func I2C(a interface{}) *Component {
	return a.(*Component)
}

// I2E - convert interface{} to *Element
func I2E(a interface{}) *E {
	return a.(*Element)
}

// IsComponent return true if interface.(type) == *Component
func IsComponent(c interface{}) bool {
	_, ok := c.(*Component)
	return ok
}

// IsElement return true if interface.(type) == *Element
func IsElement(c interface{}) bool {
	_, ok := c.(*Element)
	return ok
}

// IsString return true if interface is string
func IsString(c interface{}) bool {
	_, ok := c.(string)
	return ok
}

// RemoveStrings remove all strings from []interface{}
func RemoveStrings(arr []interface{}) []interface{} {
	var out []interface{}
	for _, el := range arr {
		if IsString(el) {
			continue
		}
		out = append(out, el)
	}
	return out
}
