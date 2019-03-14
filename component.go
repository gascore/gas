package gas

import (
	"fmt"
	"github.com/Sinicablyat/gas"
	"github.com/frankenbeanies/uuid4"
	"github.com/pkg/errors"
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

type GetComponentChildes func(*Component) []interface{}

// GetChildes -- function returning component childes
// In function parameter sends `this` component and you can get it data from this parameter
//
// Component childes can BE:
//
// 1. String (or tag_value)
//
// 2. Another component
type GetChildes func(*Component) []interface{}

// Bind - dynamic component attribute (analog for vue `v-bind:`).
type Bind func() string

// Directives struct storing component if-directive
type Directives struct {
	If    func(*Component) bool
	Else  bool
	// (If != nil) + (Else) = else-if
	Show  func(*Component) bool

	For   ForDirective
	Model ModelDirective
	HTML  HTMLDirective
}

// ModelDirective struct for Model directive
type ModelDirective struct {
	Data      string
	Component *Component
	Deep []ModelDirectiveDeepData
}

// ModelDirectiveDeepData
type ModelDirectiveDeepData struct {
	Data interface{}
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

	Get(string) Object
	Set(string, interface{})
	GetString(string) string
	GetBool(string) bool
	GetInt(string) int

	Call(string, ...interface{}) Object

	Raw() interface{}
}

// Watcher -- function triggering after component data changed
type Watcher func(*Component, interface{}, interface{}) error // (this, new, old)

// Component -- basic component struct
type Component struct {
	Data      map[string]interface{}
	Watchers  map[string]Watcher
	Methods   map[string]Method
	Computeds map[string]Computed

	Hooks Hooks // lifecycle hooks

	Handlers      map[string]Handler // events handlers: onClick, onHover
	Binds         map[string]Bind    // dynamic attributes
	RenderedBinds map[string]string  // store binds for changed func

	Directives Directives

	Childes  GetChildes
	RChildes []interface{} // rendered childes

	Tag   string
	Attrs map[string]string

	UUID string

	isElement bool // childes don't have parent context
	Parent    *Component

	Ref string

	RefsAllowed bool           // if true component can have Refs
	Refs map[string]*Component // childes have g-ref attribute. Only for Component.isElement == false

	BE BackEnd
}

// Aliases
type C = Component
type G = Gas

var NC = NewComponent
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

			childC.BE = component.BE
			childC.Parent = component

			if childC.Directives.Else {
				if !lastIfIsFresh {
					this.WarnError(errors.New("invalid else or else-if directive: no if (else-if) directive before"))
					continue
				}

				if lastIfValue {
					lastIfValue   = false
					lastIfIsFresh = childC.Directives.If != nil
					continue
				}
			}

			if childC.Directives.If != nil {
				ifValue := childC.Directives.If(childC)

				lastIfValue   = ifValue
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

func UnSpliceBody(body []interface{}) []interface{} {
	var arr []interface{}
	for _, el := range body {
		switch el.(type) {
		case []interface{}:
			for _, c := range el.([]interface{}) {
				arr = append(arr, c)
			}
			continue
		case []*gas.C:
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

	return NewForByData(dataForList, this, renderer)
}

func NewForByData(dataForList []interface{}, this *Component, renderer func(int, interface{}) interface{}) []interface{} {
	var items []interface{}
	for i, el := range dataForList {
		item := renderer(i, el)

		if IsComponent(item) {
			I2C(item).Directives.For = ForDirective{isItem: true, itemValueI: i, itemValueVal: el}

			if I2C(item).Attrs == nil {
				I2C(item).Attrs = make(map[string]string)
			}
			I2C(item).Attrs["data-for-i"] = fmt.Sprintf("%d", i)
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
	if !c.Directives.For.isItem {
		return false, 0, nil
	}

	return true, c.Directives.For.itemValueI, c.Directives.For.itemValueVal
}

// Element return *dom.Element by component
func (c *Component) Element() interface{} {
	_el := c.BE.GetElement(c)
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
	} else {
		return c.Parent.ParentComponent()
	}
}

// ParentWthAllowedRefs return first *true component* in component parents tree with allowed refs
func (c *Component) ParentWthAllowedRefs() *Component {
	// if c.Parent == nil => c - is the root (gas.App) component
	if c.Parent == nil || (!c.Parent.isElement && c.Parent.RefsAllowed) {
		return c.Parent
	} else {
		return c.Parent.ParentWthAllowedRefs()
	}
}

// GetElementUnsafely return *dom.Element by component without warning
func (c *Component) GetElementUnsafely() interface{} {
	return c.BE.GetElement(c)
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

// IsComponent return true if interface is array of interfaces
func IsChildesArr(c interface{}) bool {
	_, ok := c.([]interface{})
	return ok
}

// IsComponent return true if interface is string
func IsString(c interface{}) bool {
	_, ok := c.(string)
	return ok
}
