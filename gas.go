package gas

// Context -- in context component send c.Data and c.Props to method
type Context interface {}

// Method -- struct for Component methods
type Method func(Context)

// Component -- basic component struct
type Component struct {
	Data     map[string]interface{}
	Props    map[string]interface{}
	Methods  map[string]Method
	Childes  []Component

	Tag 	 string
	Classes  []string
	ID       string
	Attrs    map[string]string

	// Other stuff
}

// Gas -- struct for main application component
type Gas struct {
	App Component

	// Other stuff
}


// NewComponent create new component
func NewComponent(data, props *map[string]interface{}) Component {
	// Some stuff here, but now:
	return Component{
		Data: *data,
		Props: *props,
	}
}

// AddInfo add tag, id, classes and attributes to Component
func (c *Component) AddInfo(tag string, id string, classes *[]string, attrs *map[string]string) Component {
	c.Tag = tag
	c.ID = id
	c.Classes = *classes
	c.Attrs = *attrs

	return *c
}

// AddChildes add childes to Component
func (c *Component) AddChildes(childes ...Component) Component {
	c.Childes = childes

	return *c
}

// New create new gas application
//
// Here i will store explanations of incomprehensible things:
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
func New(components ...Component) (Gas, error) {
	mainComponent := NewComponent(nil, nil).AddInfo("wrap", "main", nil, nil).AddChildes(components...)
	gas := Gas{App: mainComponent}

	return gas, nil
}


func must(err error) {
	if err != nil {
		panic(err)
	}
}