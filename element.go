package gas

import (
	"fmt"
	"strings"

	"github.com/frankenbeanies/uuid4"
)

// Element stucture for basic strucutre nodes (html elements, etc)
type Element struct {
	UUID  string
	Tag   string
	Attrs map[string]string

	Handlers      map[string]Handler // events handlers: onClick, onHover
	Binds         map[string]Bind    // dynamic attributes
	RenderedBinds map[string]string  // store binds for changed func

	Watcher string
	HTML    HTMLDirective

	Childes  GetChildes
	RChildes []interface{} // rendered childes

	Parent  *Element // if element is root, parent is nil
	RefName string

	Component *Component // can be nil

	RC *RenderCore
}

// NewElement create new element
func NewElement(element *Element, childes ...interface{}) *Element {
	if element.Tag == "" {
		element.Tag = "div"
	} else {
		element.Tag = strings.ToLower(element.Tag)
	}

	element.UUID = uuid4.New().String()

	element.Childes = func() []interface{} {
		return childes
	}

	return element
}

// GetChildes function returning component/element childes
type GetChildes func() []interface{}

// Bind dynamic component attribute (analog for vue `v-bind:`).
type Bind func() string

// HTMLDirective struct for HTML Directive - storing render function and pre rendered render
type HTMLDirective struct {
	Render func() string

	Rendered string // here storing rendered html for Update functions
}

// Handler -- handler exec function when event trigger
type Handler func(Object)

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
