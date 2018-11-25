package gas

import (
	"errors"
	"github.com/dennwc/dom"
)

var (
	NilData    map[string]interface{}
	NilClasses []string
	NilAttrs   map[string]string
)

// SendComponents check components types (string or GetChildes) and return GetChildes for component.AddChildes
func SendComponents(components []interface{}) GetChildes {
	// Check components types
	for _, el := range components {
		switch el.(type) {
		case
			string,
			*Component:
			break
		default:
			panic(errors.New("invalid component type (not string or Component)"))
		}
	}

	return func(Component) []interface{} {
		return components
	}
}

// Context -- in context component send c.Data and c.Props to method
type Context interface {}

// Method -- struct for Component methods
type Method func(Context)

// GetChildes -- struct for Components Childes
// In function parameter sends `this` component and you can access component data from this parameter
//
// Component childes can be string (or tag value) like: <h1>this text is body</h1>, and another Component
type GetChildes func(Component) []interface{}

// Component -- basic component struct
type Component struct {
	Data     map[string]interface{}
	Props    map[string]interface{}
	Methods  map[string]Method
	Childes  GetChildes

	Tag 	 string
	Classes  []string
	ID       string
	Attrs    map[string]string

	// Other stuff
}

// Gas -- struct for main application component
type Gas struct {
	App Component
	StartPoint string // html element id where application wiil store

	// Other stuff
}


// NewComponent create new component
func NewComponent(data, props map[string]interface{}) *Component {
	// Some stuff here, but now:
	return &Component{
		Data: data,
		Props: props,
	}
}

// AddInfo add tag, id, classes and attributes to Component
func (c *Component) AddInfo(tag string, id string, classes []string, attrs map[string]string) *Component {
	c.Tag = tag
	c.ID = id
	c.Classes = classes
	c.Attrs = attrs

	return c
}

// AddChildes add childes to Component
//
// Component childes (tag value or Component) for works often requires component.Data(.Props) or .Methods
// and they (childes) can take this value from send in function Component
func (c *Component) AddChildes(getChildesFunction GetChildes) *Component {
	c.Childes = getChildesFunction

	return c
}

// New create new gas application
//
//
// All Add* methods return current component, if they don't we need to pre create all components before run NewComponent
// it would looks like:
//
// `
//		c1c1 := NewComponent(...)
//		c1c2 := NewComponent(...)
//
// 		c1 := NewComponent(...).AddChildes(c1c1, c1c2)
//		c2 := NewComponent(...)
//		c3 := NewComponent(...)
//
//		component := NewComponent(...).Add*(...).AddChildes(c1, c2, c3)
// ` -- seems little ridiculous
func New(startPoint string, components ...interface{}) (Gas, error) {
	mainComponent := NewComponent(NilData, NilData).AddInfo("wrap", "main", NilClasses, NilAttrs).AddChildes(SendComponents(components))

	gas := Gas{App: *mainComponent}

	el := dom.GetDocument().GetElementById(startPoint)
	if el == nil {
		return Gas{},errors.New("invalid start point")
	}

	return gas, nil
}


//func must(err error) {
//	if err != nil {
//		panic(err)
//	}
//}