package gas

import (
	"fmt"
	"strings"

	"github.com/frankenbeanies/uuid4"
)

// Element stucture for basic strucutre nodes (html elements, etc)
type Element struct {
	UUID      string
	IsPointer bool // by default element isn't pointer

	Tag      string
	Attrs    func() Map
	RAttrs   Map  // last rendered Attrs
	Handlers map[string]Handler // events handlers: onClick, onHover

	HTML HTMLDirective

	getChildes GetChildes
	Childes    []interface{}
	OldChildes []interface{}

	Parent  *Element // if element is root, parent is nil
	RefName string

	Component *Component // can be nil

	RC *RenderCore
}

// NewElement create new element
func NewElement(el *Element, childes ...interface{}) *Element {
	if el.Tag == "" {
		el.Tag = "div"
	} else {
		el.Tag = strings.ToLower(el.Tag)
	}

	if len(el.UUID) == 0 {
		el.UUID = uuid4.New().String()
	}

	el.getChildes = func() []interface{} {
		return childes
	}

	return el
}

// GetChildes function returning component/element childes
type GetChildes func() []interface{}

// Bind dynamic component attribute
type Bind func() string

// HTMLDirective struct for HTML Directive - storing render function and pre rendered render
type HTMLDirective struct {
	Render func() string

	Rendered string // here storing rendered html for Update functions
}

// Handler function for triggering event
type Handler func(Event)

// Object wrapper for js.Value
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

// Event wrapper for dom.Event
type Event interface {
	Object
	Value() string
}

// BEElement return element in backend implementation
func (e *Element) BEElement() interface{} {
	_el := e.RC.BE.GetElement(e)
	if _el == nil {
		e.RC.BE.ConsoleError(fmt.Sprintf("component Element: %s, returning nil", e.UUID))
		return nil
	}

	return _el
}

// GetElementUnsafely return *dom.Element by component without warning
func (e *Element) GetElementUnsafely() interface{} {
	return e.RC.BE.GetElement(e)
}

// ParentComponent return first component in element parents tree
func (e *Element) ParentComponent() *Element {
	parent := e.Parent
	if parent.Parent == nil || parent.Component != nil {
		return parent
	}

	return parent.ParentComponent()
}
