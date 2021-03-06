package gas

// Template function returning elements
type Template func(...interface{}) []interface{}

// External structure for passing values to external components
type External struct {
	Body      []interface{}
	Slots     map[string]interface{}
	Templates map[string]Template
	Attrs     func() Map
}

type DynamicElement func(External) *Element

type DynamicComponent func(External) *Component
