package gas

// Component logic node
type Component struct {
	Root interface {
		Render() *Element
	}

	Element *Element // last-time render by root element from
	Hooks   Hooks

	RefsAllowed bool
	Refs        map[string]*Element

	NotPointer bool // by default component is pointer

	RC *RenderCore
}

// RenderElement create element for elements tree
func (c *Component) RenderElement() *Element {
	el := c.Root.Render()

	el.RC = c.RC
	el.Component = c
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
	C       *Component
	Element *Element
}

func (root *EmptyRoot) Render() *Element {
	return root.Element
}
