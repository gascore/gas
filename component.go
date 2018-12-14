package gas

import (
	"fmt"
	"github.com/Sinicablyat/dom"
	"github.com/frankenbeanies/uuid4"
	"strings"
)

var (
	// NilData Nil value for Component.Data
	NilData map[string]interface{}
	// NilWatchers Nil value for Component.Watchers
	NilWatchers map[string]Watcher
	// NilAttrs Nil value for Component.Attrs
	NilAttrs map[string]string
	// NilBinds Nil value for Component.Binds
	NilBinds map[string]Bind
	// NilHooks Nil value for Component.Hooks
	NilHooks Hooks
	// NilHandlers Nil value for Component.Handlers
	NilHandlers map[string]Handler
	// NilMethods Nil value for Component.Methods
	NilMethods map[string]Method
	// NilComputeds Nil value for Component.Computeds
	NilComputeds map[string]Computed
	// NilDirectives Nil value for Component.Directives
	NilDirectives Directives
)

// Context - in context component send c.Data and c.Props to method
type Context interface{}

// Method - struct for Component methods
type Method func(*Component, ...interface{}) error

// Computed - struct for Component computed values
type Computed func(*Component, ...interface{}) (interface{}, error)

// GetComponent returns component child
type GetComponent func(*Component) interface{}

// Hooks component lifecycle hooks
type Hooks struct {
	Created 	 Hook // have auto-update
	BeforeCreate Hook
	Destroyed 	 Hook
}

// Hook - lifecycle hook
type Hook func(*Component) error

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
type Bind func(*Component) string

// Catcher catching sub component $emit and doing his business  (analogue for vue `v-on:`). 
type Catcher func(*Component)

// Directives struct storing component if-directive
type Directives struct {
	If func(*Component) bool
	For ForDirective
	Model ModelDirective
	HTML HTMLDirective
}

// ModelDirective struct for Model directive
type ModelDirective struct {
	Data string
	Component *Component
}

// ForDirective struct for For Directive (needful because `for` want name and render function)
type ForDirective struct {
	Data string
	Render func(int, interface{}, *Component) []GetComponent
}

// HTMLDirective struct for HTML Directive - storing render function and pre rendered render
type HTMLDirective struct {
	Render func(*Component) string

	Rendered string // here storing rendered html for update functions
}

// Handler -- handler exec function when event trigger
type Handler func(*Component, dom.Event)

// Watcher -- function triggering after component data changed
type Watcher func(*Component, interface{}, interface{})error // (this, new, old)


// Component -- basic component struct
type Component struct {
	Data  map[string]interface{}
	Watchers map[string]Watcher

	Methods    	 map[string]Method
	Computeds    map[string]Computed

	Hooks    Hooks // lifecycle hooks
	Catchers map[string]Catcher // catch child components $emit

	Handlers      map[string]Handler // events handlers: onClick, onHover
	Binds      	  map[string]Bind    // dynamic attributes
	renderedBinds map[string]string // store binds for changed func

	Directives 	 Directives

	Childes GetChildes
	RChildes []interface{} // rendered childes

	Tag   string
	Attrs map[string]string

	UUID string

	ParentC *Component
}

func NewComponent(component *Component, childes ...GetComponent) *Component {
	component.Tag = strings.ToLower(component.Tag)
	component.UUID = uuid4.New().String()

	component.Childes = func(this *Component) []interface{} {
		var compiled []interface{}
		for _, el := range childes {
			child := el(this)

			if isComponent(child) {
				childC := I2C(child)
				if childC.Directives.If != nil && !childC.Directives.If(childC) {
					continue
				}

				// if for.Data doesn't exist, but render exist - it's a user problem
				if childC.Directives.For.Render != nil {
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

						c := &Component{
							ParentC: childC.ParentC,
							Data: childC.Data,
							Watchers: component.Watchers,
							Methods: childC.Methods,
							Computeds: component.Computeds,
							Directives: clearedDirective,
							Binds: childC.Binds,
							Hooks: component.Hooks,
							Handlers: childC.Handlers,
							Tag: childC.Tag,
							Attrs: childC.Attrs,
						}

						oneOfComponents := NewComponent(c, renderer(i, el, this)...)
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
func (c *Component) GetElement() *dom.Element {
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
