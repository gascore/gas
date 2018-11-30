package gas

import (
	"errors"
	"fmt"
	"github.com/Sinicablyat/dom"
	"github.com/frankenbeanies/uuid4"
)

var (
	// NilParentComponent Nil value for Component.ParentC
	NilParentComponent *Component
	// NilData Nil value for Component.Data and .Props
	NilData map[string]interface{}
	// NilAttrs Nil value for Component.Attrs
	NilAttrs map[string]string
	// NilBinds Nil value for Component.Binds
	NilBinds map[string]Bind
	// NilHandlers Nil value for Component.Handlers
	NilHandlers map[string]Handler
	// NilMethods Nil value for Component.Methods
	NilMethods map[string]Method
)

// Context - in context component send c.Data and c.Props to method
type Context interface{}

// Method - struct for Component methods
type Method func(Context)

// GetComponent returns component child
type GetComponent func(Component) interface{}

// GetChildes -- function returning component childes
// In function parameter sends `this` component and you can get it data from this parameter
//
// Component childes can be:
//
// 1. String (or tag_value)
//
// 2. Another component
type GetChildes func(Component) []interface{}

// Bind -- bind catching sub component $emit and doing his business.
// It's analogue for vue `v-bind:`
// Like: `gBind:id="c.GetDataByString("iterator") + 1024"``
type Bind func(Component) string

// Handler -- handler exec function when event trigger
type Handler func(Component, dom.Event)

// Component -- basic component struct
type Component struct {
	Data  map[string]interface{}
	Props map[string]interface{}

	Methods  map[string]Method
	Handlers map[string]Handler // events handlers: onClick, onHover
	Binds    map[string]Bind   // catch sub components $emit

	Childes GetChildes

	Tag   string
	Attrs map[string]string

	UUID string

	ParentC *Component
}

// NewComponent create new component
func NewComponent(pC *Component, data map[string]interface{}, methods map[string]Method, binds map[string]Bind, handlers map[string]Handler, tag string, attrs map[string]string, childes ...GetComponent) *Component {
	// Some stuff here, but now:
	component := &Component{
		Data:  data,

		Methods: methods,
		Handlers: handlers,
		Binds: binds,

		Tag:   tag,
		Attrs: attrs,

		UUID: uuid4.New().String(),

		ParentC: pC,
	}

	component.Childes = func(this Component) []interface{} {
		var compiled []interface{}
		for _, el := range childes {
			compiled = append(compiled, el(this))
		}

		return compiled
	}

	return component
}

// CreateComponent render component
func CreateComponent(node interface{}) (*dom.Element, error) {
	switch node.(type) {
	case string:
		nodeS := node.(string)
		_node := dom.NewElement("span") // dennwc/dom doesn't support textNode

		if _node == nil {
			return nil, errors.New("cannot create textNode")
		}

		_node.SetTextContent(nodeS)

		return _node, nil
	case *Component:
		component := node.(*Component)

		_node := dom.NewElement(component.Tag)
		if _node == nil {
			return nil, errors.New("cannot create component")
		}

		_node.SetAttribute("data-i", component.UUID) // set data-i for accept element from component methods

		for attrName, attrBody := range component.Attrs {
			//if attrIsValid(attrName, component.Tag) {} // check if attribute is valid for this tag
			_node.SetAttribute(attrName, attrBody)
		}

		for handlerName, handlerBody := range component.Handlers {
			//if handlerIsValid(handlerName, component.Tag) {} // check if handler is valid for this tag
			_node.AddEventListener(handlerName, func(e dom.Event) {
				handlerBody(*component, e)
			})
		}

		for _, el := range component.Childes(*component) {
			_child, err := CreateComponent(el)
			if err != nil {
				return nil, err
			}

			_node.AppendChild(_child)
		}

		return _node, nil
	case nil:
		return CreateComponent("nil")
	default:
		err := fmt.Errorf("invalid component type: %T", node)
		return nil, err
	}
}

// UpdateComponent trying to update component
func UpdateComponent(_parent *dom.Element, new interface{}, old interface{}, index int) error {
	// if component has created
	if old == nil {
		_new, err := CreateComponent(new)
		if err != nil {
			return err
		}

		_parent.AppendChild(_new)

		return nil
	}

	_childes := _parent.ChildNodes()
	if _childes == nil {
		return errors.New("_parent doesn't have childes")
	}
	if len(_childes) <= index { // check index on valid
		return errors.New("invalid index")
	}
	_el := _childes[index]

	// if component has deleted
	if new == nil {
		_parent.RemoveChild(_el)

		return nil
	}

	// if component has changed
	isChanged, err := changed(new, old)
	if err != nil {
		return err
	}
	if isChanged {
		_new, err := CreateComponent(new)
		if err != nil {
			return err
		}

		_old := _el

		_parent.ReplaceChild(_new, _old)
	}

	// if component childes have changed
	// check that `new` and `old` is components
	if isComponent(new) && isComponent(old) {
		newChildes := new.(*Component).Childes(*new.(*Component)) // new.Childes(new)
		oldChildes := old.(*Component).Childes(*old.(*Component)) // old.Childes(old)

		err := UpdateComponentChildes(_el, newChildes, oldChildes)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateComponentChildes update component childes by new and old childes
func UpdateComponentChildes(_el *dom.Element, newChildes, oldChildes []interface{}) error {
	for i := 0; i < len(newChildes) || i < len(oldChildes); i++ {
		var elFromNew interface{}
		if len(newChildes) >= i {
			elFromNew = newChildes[i]
		}

		var elFromOld interface{}
		if len(oldChildes) >= i {
			elFromOld = oldChildes[i]
		}

		err := UpdateComponent(_el, elFromNew, elFromOld, i)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetData return data field by query string
func (c *Component) GetData(query string) interface{} {
	// There will be callbacks, events, e.t.c.
	data := c.Data[query]
	if data == nil {
		dom.ConsoleError(fmt.Sprintf(`"%s"trying to accept nil data`, c.Tag))
	}

	return data
}

// SetData set data field and update component (after changes)
func (c *Component) SetData(query string, value interface{}) error {
	oldChildes := c.Childes(*c)
	c.Data[query] = value

	_c := c.GetElement()
	dom.ConsoleLog(_c)

	newChildes := c.Childes(*c)

	err := UpdateComponentChildes(_c, newChildes, oldChildes)
	if err != nil {
		dom.ConsoleError(err)
		return err
	}

	return nil
}

// GetElement return *dom.Element by component structure
func (c Component) GetElement() *dom.Element {
	return dom.Doc.QuerySelector(fmt.Sprintf("[data-i='%s']", c.UUID)) // select element by data-i attribute
}


func isComponent(c interface{}) bool {
	_, ok := c.(*Component)
	return ok
}

func changed(new, old interface{}) (bool, error) {
	newType := fmt.Sprintf("%T", new)
	oldType := fmt.Sprintf("%T", old)

	if newType != oldType {
		return true, nil
	}

	if newType == "string" && newType == oldType {
		newString, ok := new.(string)
		if !ok {
			return false, errors.New("invalid `new`")
		}

		oldString, ok := old.(string)
		if !ok {
			return false, errors.New("invalid `old`")
		}

		return newString != oldString, nil
	}

	if isComponent(new) && isComponent(old) {
		newComponent, ok := new.(*Component)
		if !ok {
			return false, errors.New("invalid `new`")
		}

		oldComponent, ok := old.(*Component)
		if !ok {
			return false, errors.New("invalid `old`")
		}

		return newComponent != oldComponent, nil
	}

	return false, fmt.Errorf("changed: invalid `new` or `old`. types: %T, %T", new, old)
}
