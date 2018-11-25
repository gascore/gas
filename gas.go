package gas

import (
	"errors"
	"github.com/dennwc/dom"
)

var (
	NilData    map[string]interface{}
	NilClasses []string
	NilAttrs   map[string]string
	NilBinds   map[string]Bind
)

// Context -- in context component send c.Data and c.Props to method
type Context interface {}

// Method -- struct for Component methods
type Method func(Context)

// GetComponent -- function returning component child
type GetComponent func(Component) interface{}

// GetChildes -- function returning component childes
// In function parameter sends `this` component and you can get component data from this parameter
//
// Component childes can be :
// 1. String (or tag_value) like: <h1>this text is body</h1>
// 2. Another component like: <h1> <GreetingText/> </h1>
type GetChildes func(Component) []interface{}

// Bind -- struct for storage component *computed* attributes.
// It's analouge for vue `v-bind:`
// Like: `gBind:id="c.GetDataByString("iterator") + 1024"`` 
type Bind func(Component) string

// Component -- basic component struct
type Component struct {
	Data      map[string]interface{}
	Props     map[string]interface{}
	
	CallBacks map[string]Method // events handlers: onClick, onHover
	Methods   map[string]Method // user functions can call from component childes
	
	Childes   GetChildes

	Tag 	  string
	ID        string
	Classes   []string
	Attrs	  map[string]string

	Binds     map[string]Bind

	// Other stuff
}

// Gas -- struct for main application component
type Gas struct {
	App Component
	StartPoint string // html element id where application wiil store

	// Other stuff
}


// NewComponent create new component
func NewComponent(data, props map[string]interface{}, tag string, id string, classes []string, attrs map[string]string) *Component {
	// Some stuff here, but now:
	return &Component{
		Data: data,
		Props: props,

		Tag: tag,
		ID: id,
		Classes: classes,
	}
}

// AddBinds add Binds to Component
func (c *Component) AddBinds(binds map[string]Bind) *Component {
	c.Binds = binds

	return c
}

// AddChildes add childes to Component
//
// Component childes (tag value or Component) for works often requires component.Data(.Props) or .Methods
// and they (childes) can take this value from send in function Component
func (c *Component) AddChildes(childes ...GetComponent) *Component {
	c.Childes = func(this Component) []interface{} {
		var compiled []interface{}
		for _, el := range childes {
			compiled = append(compiled, el(this))
		}

		return compiled
	}

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
func New(startPoint string, components ...GetComponent) (Gas, error) {
	mainComponent := NewComponent(NilData, NilData, "wrap", "main", NilClasses, NilAttrs).AddBinds(NilBinds).AddChildes(components...)

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