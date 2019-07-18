package gas

import "github.com/frankenbeanies/uuid4"

// Component logic node
type Component struct {
	Root interface {
		Render() []interface{}
	}

	Element            *Element
	ElementIsImportant bool

	Hooks Hooks

	RefsAllowed bool
	Refs        map[string]*Element

	Watchers map[string]Watcher

	NotPointer bool // by default component is pointer

	RC *RenderCore
}

// Watcher function called when input event triggering
type Watcher func(val interface{}, e Object) (string, error)

// Init initialize component: create element and other stuff
func (c *Component) Init() *Element {
	var el *Element
	if c.Element == nil {
		el = &Element{
			Tag:       "div",
			UUID:      uuid4.New().String(),
			Component: c,
			RC:        c.RC,
		}
	} else {
		el = c.Element
		el.Component = c
		el.RC = c.RC

		if len(el.UUID) == 0 {
			el.UUID = uuid4.New().String()
		}

		if len(el.Tag) == 0 {
			el.Tag = "div"
		}
	}

	el.getChildes = func() []interface{} {
		return c.Root.Render()
	}
	el.IsPointer = !c.NotPointer

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

// EmptyRoot root for component only rendering one element
type EmptyRoot struct {
	C *C
	Element *Element
}

func (root *EmptyRoot) Render() []interface{} {
	return CL(root.Element)
}
