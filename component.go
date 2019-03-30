package gas

import (
	"fmt"
	"strings"

	"github.com/frankenbeanies/uuid4"
	"github.com/pkg/errors"
)

// Context - in context component send c.Data and c.Props to method
type Context interface{}

// Method - struct for Component methods
type Method func(*Component, ...interface{}) (interface{}, error)

// GetComponent returns component child
type GetComponent func(*Component) interface{}

// GetComponentChildes return component childes
type GetComponentChildes func(*Component) []interface{}

// GetChildes -- function returning component childes
type GetChildes func(*Component) []interface{}

// Bind - dynamic component attribute (analog for vue `v-bind:`).
type Bind func() string

// ModelDirective struct for Model directive
type ModelDirective struct {
	Data      string
	Component *Component
	Deep      []ModelDirectiveDeepData
}

// ModelDirectiveDeepData way to deep value of data field need to update by model directive
type ModelDirectiveDeepData struct {
	Data     interface{}
	Brackets bool
}

// ForDirective struct for For Directive (needful because `for` want name and render function)
type ForDirective struct {
	isItem       bool
	itemValueI   int
	itemValueVal interface{}
}

// HTMLDirective struct for HTML Directive - storing render function and pre rendered render
type HTMLDirective struct {
	Render func(*Component) string

	Rendered string // here storing rendered html for ForceUpdate functions
}

// Handler -- handler exec function when event trigger
type Handler func(*Component, Object)

// Object 'united' dom.Event
type Object interface {
	String() string
	Int() int
	Float() float64

	Get(string) Object
	Set(string, interface{})
	GetString(string) string
	GetBool(string) bool
	GetInt(string) int

	Call(string, ...interface{}) Object

	Raw() interface{}
}

// Watcher -- function triggering after component Data changed
type Watcher func(this *Component, new interface{}, old interface{}) error // (this, new, old)

// Component -- basic component struct
type Component struct {
	Data     map[string]interface{}
	Watchers map[string]Watcher
	Methods  map[string]Method

	Hooks Hooks // lifecycle hooks

	Handlers      map[string]Handler // events handlers: onClick, onHover
	Binds         map[string]Bind    // dynamic attributes
	RenderedBinds map[string]string  // store binds for changed func

	/* directives */
	If   func(*Component) bool
	Else bool
	// (If != nil) + (Else) = else-if
	Show  func(*Component) bool
	For   ForDirective
	Model ModelDirective
	HTML  HTMLDirective

	Childes  GetChildes
	RChildes []interface{} // rendered childes

	Tag   string
	Attrs map[string]string

	UUID string

	isElement bool // childes don't have parent context
	Parent    *Component

	Ref         string
	RefsAllowed bool                  // if true component can have Refs
	Refs        map[string]*Component // childes have g-ref attribute. Only for Component.isElement == false

	RC *RenderCore
}

// C alias for Component
type C = Component

// G alias for Gas
type G = Gas

// NC alias for NewComponent
var NC = NewComponent

// NE alias for NewBasicComponent (NewElement)
var NE = NewBasicComponent

// NewComponent create new component
func NewComponent(component *Component, getChildes GetComponentChildes) *Component {
	if component.Tag == "" {
		component.Tag = "div"
	} else {
		component.Tag = strings.ToLower(component.Tag)
	}

	if component.UUID == "" {
		component.UUID = uuid4.New().String()
	}

	component.Childes = func(this *Component) []interface{} {
		var lastIfValue, lastIfIsFresh bool
		var compiled []interface{}
		for _, child := range UnSpliceBody(getChildes(component)) {
			if !IsComponent(child) {
				compiled = append(compiled, child)
				continue
			}

			childC := I2C(child)

			childC.RC = component.RC
			childC.Parent = component

			if childC.Else {
				if !lastIfIsFresh {
					this.WarnError(errors.New("invalid else or else-if directive: no if (else-if) directive before"))
					continue
				}

				if lastIfValue {
					lastIfValue = false
					lastIfIsFresh = childC.If != nil
					continue
				}
			}

			if childC.If != nil {
				ifValue := childC.If(childC)

				lastIfValue = ifValue
				lastIfIsFresh = true

				if !ifValue {
					continue
				}
			}

			compiled = append(compiled, child)
		}

		return compiled
	}

	return component
}

// NewBasicComponent create new component without *this* context
func NewBasicComponent(component *Component, childes ...interface{}) *Component {
	component.isElement = true
	return NewComponent(component, func(this *Component) []interface{} {
		return childes
	})
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
		case []*C:
			for _, c := range el.([]interface{}) {
				arr = append(arr, c)
			}
			continue
		default:
			arr = append(arr, el)
		}
	}

	return arr
}

// NewFor create new FOR directive
func NewFor(data string, this *Component, renderer func(int, interface{}) interface{}) []interface{} {
	dataForList, ok := this.Data[data].([]interface{})
	if !ok {
		this.WarnError(fmt.Errorf("invalid FOR directive in component %s", this.UUID))
		return nil
	}

	return NewForByData(dataForList, renderer)
}

// NewForByData create new FOR directive by []interface{}
func NewForByData(dataForList []interface{}, renderer func(int, interface{}) interface{}) []interface{} {
	var items []interface{}
	for i, el := range dataForList {
		item := renderer(i, el)

		if IsComponent(item) {
			I2C(item).For = ForDirective{isItem: true, itemValueI: i, itemValueVal: el}

			if I2C(item).Attrs == nil {
				I2C(item).Attrs = make(map[string]string)
			}
			I2C(item).Attrs["Data-for-i"] = fmt.Sprintf("%d", i)
		}

		items = append(items, item)
	}

	return items
}

// IsElement return Component.isElement.
// isElement isn't public because it's useless for applications developing but for helpful libraries
func (c *Component) IsElement() bool {
	return c.isElement
}

// ForItemInfo return info about FOR directive
func (c *Component) ForItemInfo() (isItem bool, i int, val interface{}) {
	if !c.For.isItem {
		return false, 0, nil
	}

	return true, c.For.itemValueI, c.For.itemValueVal
}

// Element return *dom.Element by component
func (c *Component) Element() interface{} {
	_el := c.RC.BE.GetElement(c)
	if _el == nil {
		c.WarnError(fmt.Errorf("component Element: %s, returning nil", c.UUID))
		return nil
	}

	return _el
}

// ParentComponent return first *true component* in component parents tree
func (c *Component) ParentComponent() *Component {
	// if c.Parent == nil => c - is the root (gas.App) component
	if c.Parent == nil || !c.Parent.isElement {
		return c.Parent
	}

	return c.Parent.ParentComponent()
}

// ParentWithAllowedRefs return first *true component* in component parents tree with allowed refs
func (c *Component) ParentWithAllowedRefs() *Component {
	// if c.Parent == nil => c - is the root (gas.App) component
	if c.Parent == nil || (!c.Parent.isElement && c.Parent.RefsAllowed) {
		return c.Parent
	}

	return c.Parent.ParentWithAllowedRefs()
}

// GetElementUnsafely return *dom.Element by component without warning
func (c *Component) GetElementUnsafely() interface{} {
	return c.RC.BE.GetElement(c)
}

// I2C - convert interface{} to *Component
func I2C(a interface{}) *Component {
	return a.(*Component)
}

// IsComponent return true if interface is *Component
func IsComponent(c interface{}) bool {
	_, ok := c.(*Component)
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
