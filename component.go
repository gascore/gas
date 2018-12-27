package gas

import (
	"fmt"
	"github.com/frankenbeanies/uuid4"
	"strings"
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
	Created 	 Hook // have auto-Update
	BeforeCreate Hook
	Destroyed 	 Hook
	BeforeUpdate Hook
	Updated		 Hook
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

// Directives struct storing component if-directive
type Directives struct {
	If func(*Component) bool
	Show func(*Component) bool
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
	isItem bool
	itemValueI int
	itemValueVal interface{}
}

// HTMLDirective struct for HTML Directive - storing render function and pre rendered render
type HTMLDirective struct {
	Render func(*Component) string

	Rendered string // here storing rendered html for Update functions
}

// Handler -- handler exec function when event trigger
type Handler func(*Component, HandlerEvent)

// HandlerEvent 'united' dom.Event
type HandlerEvent interface {
	Get(string) interface{}
	GetString(string) string
	GetBool(string) bool
	GetInt(string) int

	Raw() interface{}
}

// Watcher -- function triggering after component data changed
type Watcher func(*Component, interface{}, interface{})error // (this, new, old)


// Component -- basic component struct
type Component struct {
	Data  map[string]interface{}
	Watchers map[string]Watcher

	Methods    	 map[string]Method
	Computeds    map[string]Computed

	Hooks    Hooks // lifecycle hooks
	Handlers      map[string]Handler // events handlers: onClick, onHover
	Binds      	  map[string]Bind    // dynamic attributes
	RenderedBinds map[string]string // store binds for changed func

	Directives 	 Directives

	Childes GetChildes
	RChildes []interface{} // rendered childes

	Tag   string
	Attrs map[string]string

	UUID string

	ParentC *Component
}

func NewComponent(component *Component, childes ...GetComponent) *Component {
	if component.Tag == "" {
		component.Tag = "div"
	}

	component.Tag = strings.ToLower(component.Tag)

	component.Childes = func(this *Component) []interface{} {
		var compiled []interface{}
		for _, el := range childes {
			child := el(this)

			if IsComponent(child) {
				childC := I2C(child)
				if childC.Directives.If != nil && !childC.Directives.If(childC) {
					continue
				}
			} else if IsChildesArr(child) {
				compiled = append(compiled, child.([]interface{})...)
				continue
			}

			compiled = append(compiled, child)
		}

		return compiled
	}

	component.UUID = uuid4.New().String()

	return component
}

func NewFor(data string, this *Component, renderer func(int, interface{}) interface{}) []interface{} {
	dataForList, ok := this.Data[data].([]interface{})
	if !ok {
		WarnError(fmt.Errorf("invalid FOR directive in component %s", this.UUID))
		return nil
	}

	var items []interface{}
	for i, el := range dataForList {
		item := renderer(i, el)

		if IsComponent(item) {
			itemC := I2C(item)
			itemC.Directives.For = ForDirective{isItem: true, itemValueI: i, itemValueVal: el}
		}

		items = append(items, item)
	}

	return items
}

func (c *Component) ForItemInfo() (bool, int, interface{}) {
	if !c.Directives.For.isItem {
		return false, 0, nil
	}

	return true, c.Directives.For.itemValueI, c.Directives.For.itemValueVal
}

// GetElement return *dom.Element by component structure
func (c *Component) GetElement() interface{} {
	return be.GetElement(c)
}

// I2C - convert interface{} to *Component
func I2C(a interface{}) *Component {
	return a.(*Component)
}


func IsComponent(c interface{}) bool {
	_, ok := c.(*Component)
	return ok
}
func IsChildesArr(c interface{}) bool {
	_, ok := c.([]interface{})
	return ok
}
func IsString(c interface{}) bool {
	_, ok := c.(string)
	return ok
}
