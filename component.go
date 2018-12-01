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
	// NilDirectives Nil value for Component.Directives
	NilDirectives = Directives{If:NilIfDirective}
	// NilIfDirective Nil value for Directives.If
	NilIfDirective = func(c *Component) bool { return true } // Without returning component will never render
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
type GetChildes func(*Component) []interface{}

// Bind -- bind catching sub component $emit and doing his business.
// It's analogue for vue `v-bind:`
// Like: `gBind:id="c.GetDataByString("iterator") + 1024"``
type Bind func(Component) string

// Directives
type Directives struct {
	If func(*Component) bool
	For func(arr []interface{}) []Component
}

// Handler -- handler exec function when event trigger
type Handler func(Component, dom.Event)

// Component -- basic component struct
type Component struct {
	Data  map[string]interface{}
	Props map[string]interface{}

	Methods    map[string]Method
	Handlers   map[string]Handler // events handlers: onClick, onHover
	Binds      map[string]Bind   // catch sub components $emit
	Directives Directives

	Childes GetChildes

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
			child   := el(*this)

			if isComponent(child) && !I2C(child).Directives.If(I2C(child)) {
				continue
			}

			compiled = append(compiled, child)
		}

		return compiled
	}

	return component
}

// CreateComponent render component. Returns _el, err
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
		component := I2C(node)

		_node, err := CreateElement(component)
		if err != nil {
			return nil, errors.New("cannot create component")
		}

		for _, el := range component.Childes(component) {
			_child, err := CreateComponent(el)
			if err != nil {
				return nil, err
			}

			if _child != nil {
				_node.AppendChild(_child)
			}
		}

		return _node, nil
	default:
		return nil, fmt.Errorf("invalid component type: %T", node)
	}
}

// CreateElement create html element without childes
func CreateElement(c *Component) (*dom.Element, error) {
	_node := dom.NewElement(c.Tag)
	if _node == nil {
		return nil, errors.New("cannot create component")
	}

	_node.SetAttribute("data-i", c.UUID) // set data-i for accept element from component methods

	for attrName, attrBody := range c.Attrs {
		//if attrIsValid(attrName, component.Tag) {} // check if attribute is valid for this tag
		_node.SetAttribute(attrName, attrBody)
	}

	for handlerName, handlerBody := range c.Handlers {
		//if handlerIsValid(handlerName, component.Tag) {} // check if handler is valid for this tag
		_node.AddEventListener(handlerName, func(e dom.Event) {
			handlerBody(*c, e)
		})
	}

	return _node, nil
}

// UpdateComponent trying to update component
func UpdateComponent(_parent *dom.Element, new interface{}, old interface{}, index int) error {
	// if component has created
	if old == nil {
		_new, err := CreateComponent(new)
		if err != nil {
			return err
		}

		if _new != nil {
			_parent.AppendChild(_new)
		}

		return nil
	}

	_childes := _parent.ChildNodes()
	if _childes == nil { return errors.New("_parent doesn't have childes") }

	var _el *dom.Element
	if len(_childes) > index { // component was hided if childes length <= index
		_el = _childes[index]
	}

	newIsComponent := isComponent(new)
	var newC *Component
	if newIsComponent {
		newC = I2C(new)
	}

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

		if _new != nil {
			_parent.ReplaceChild(_new, _el)
		}

		return nil
	}

	// check if component childes updated
	if newIsComponent {
		_new, err := CreateComponent(new)
		if err != nil {
			return err
		}

		if _new != nil {
			_parent.ReplaceChild(_new, _el)
		}

		return nil
		// update element
		//currentChildes := _el.ChildNodes() // get childes from old _el
		//log.Println(currentChildes)
		//
		//_newEl, err := CreateElement(newC) // create new _el
		//if err != nil {
		//	return err
		//}
		//log.Println(_newEl)
		//
		//for _, _child := range currentChildes { // transfer childes from old _el to new
		//	//_newEl.AppendChild(child)
		//	dom.ConsoleLog(_child)
		//	_newEl.AppendChild(_child)
		//}

		// update childes
		newChildes := newC.Childes(newC) // new.Childes(new)
		oldChildes := I2C(old).Childes(I2C(old)) // old.Childes(old)

		err = UpdateComponentChildes(_el, newChildes, oldChildes)
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

		//if isComponent(elFromNew) && isComponent(elFromOld) {log.Println(elFromNew, elFromOld)}
		err := UpdateComponent(_el, elFromNew, elFromOld, i)
		if err != nil {
			return err
		}
	}
	return nil
}

// changed return true if node changed
func changed(new, old interface{}) (bool, error) {
	if fmt.Sprintf("%T", new) != fmt.Sprintf("%T", old) {
		return true, nil
	}

	if isString(new) {
		return new.(string) != old.(string), nil
	} else if isComponent(new) {
		return false, nil
	}

	return false, fmt.Errorf("changed: invalid `new` or `old`. types: %T, %T", new, old)
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
	oldChildes := c.Childes(c)


	c.Data[query] = value


	newChildes := c.Childes(c)
	_c := c.GetElement()

	err := UpdateComponentChildes(_c, newChildes, oldChildes)
	if err != nil {
		return err
	}

	return nil
}



// GetElement return *dom.Element by component structure
func (c Component) GetElement() *dom.Element {
	return dom.Doc.QuerySelector(fmt.Sprintf("[data-i='%s']", c.UUID)) // select element by data-i attribute
}

// I2C - convert interface{} to *Component
func I2C(a interface{}) *Component {
	return a.(*Component)
}

func isComponent(c interface{}) bool {
	_, ok := c.(*Component)
	return ok
}
func isString(c interface{}) bool {
	_, ok := c.(string)
	return ok
}
