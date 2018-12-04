package gas

import (
	"fmt"
	"github.com/Sinicablyat/dom"
	"github.com/frankenbeanies/uuid4"
)

var (
	// NilParentComponent Nil value for Component.ParentC
	NilParentComponent *Component
	// NilData Nil value for Component.Data
	NilData map[string]interface{}
	// NilAttrs Nil value for Component.Attrs
	NilAttrs map[string]string
	// NilBinds Nil value for Component.Binds
	NilBinds map[string]Bind
	// NilHandlers Nil value for Component.Handlers
	NilHandlers map[string]Handler
	// NilMethods Nil value for Component.Methods
	NilMethods map[string]Method
	// NilDirectives Nil value for Component.Directives
	NilDirectives = Directives{If:NilIfDirective}
	// NilIfDirective Nil value for Directives.If
	NilIfDirective = func(c *Component) bool { return true } // Without returning component will never render
)

// Context - in context component send c.Data and c.Props to method
type Context interface{}

// Method - struct for Component methods
type Method func(Context) interface{}

// GetComponent returns component child
type GetComponent func(*Component) interface{}

// GetChildes -- function returning component childes
// In function parameter sends `this` component and you can get it data from this parameter
//
// Component childes can be:
//
// 1. String (or tag_value)
//
// 2. Another component
type GetChildes func(*Component) []interface{}

// Bind - dynamic component attribute (analogue for vue `v-bind:`).
//
// Like: `gBind:id="c.GetDataByString("iterator") + 1024"``
type Bind func(Component) string

// Catcher catching sub component $emit and doing his business  (analogue for vue `v-on:`). 
type Catcher func(Component)

// Directives struct storing component if-directive
type Directives struct {
	If func(*Component) bool
	For ForDirective
}

// ForDirective struct for For Directive (needful because `for` want name and render function)
type ForDirective struct {
	Data string
	Render func(int, interface{}, *Component) []GetComponent
}

// Handler -- handler exec function when event trigger
type Handler func(Component, dom.Event)


// Component -- basic component struct
type Component struct {
	Data  map[string]interface{}

	Catchers map[string]Catcher // catch child components $emit
	Methods    	 map[string]Method
	Handlers     map[string]Handler // events handlers: onClick, onHover
	Binds      	 map[string]Bind    // dynamic attributes
	Directives 	 Directives

	Childes GetChildes
	RChildes []interface{} // rendered childes

	Tag   string
	Attrs map[string]string

	UUID string

	ParentC *Component
}


// NewComponent create new component
func NewComponent(pC *Component, data map[string]interface{}, methods map[string]Method, directives Directives, binds map[string]Bind, handlers map[string]Handler, tag string, attrs map[string]string, childes ...GetComponent) *Component {
	// Some stuff here, but now:
	component := &Component{
		Data:  data,

		Methods: methods,
		Handlers: handlers,
		Binds: binds,
		Directives: directives,

		Tag:   tag,
		Attrs: attrs,

		UUID: uuid4.New().String(),

		ParentC: pC,
	}

	component.Childes = func(this *Component) []interface{} {
		var compiled []interface{}
		for _, el := range childes {
			child := el(this)

			if isComponent(child) {
				childC := I2C(child)
				if !childC.Directives.If(childC) {
					continue
				}

				// if for.Data doesn't exist, but render exist - it's a user problem
				if childC.Directives.For.Data != "" {
					dataForList, ok := this.Data[childC.Directives.For.Data].([]interface{})
					if !ok {
						dom.ConsoleError(fmt.Sprintf("invalid FOR directive in component %s", childC.UUID))
						continue
					}

					clearedDirective := childC.Directives
					clearedDirective.For = ForDirective{}

					renderer := childC.Directives.For.Render
					for i, el := range dataForList {
						// recreate this component with childes from FOR, without FOR directive
						oneOfComponents := NewComponent(childC.ParentC, childC.Data, childC.Methods, clearedDirective, childC.Binds, childC.Handlers, childC.Tag, childC.Attrs, renderer(i, el, this)...)
						compiled = append(compiled, oneOfComponents)
					}

					continue
				}
			}

			compiled = append(compiled, child)
		}

		return compiled
	}

	return component
}

// GetElement return *dom.Element by component structure
func (c Component) GetElement() *dom.Element {
	return dom.Doc.QuerySelector(fmt.Sprintf("[data-i='%s']", c.UUID)) // select element by data-i attribute
}

// I2C - convert interface{} to *Component
func I2C(a interface{}) *Component {
	return a.(*Component)
}

// ToGetComponentList return array by many parameters, because it's pretty
func ToGetComponentList(childes ...GetComponent) []GetComponent {return childes}

func isComponent(c interface{}) bool {
	_, ok := c.(*Component)
	return ok
}
func isString(c interface{}) bool {
	_, ok := c.(string)
	return ok
}
