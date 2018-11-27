package gas

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dennwc/dom"
)

var (
	// NilData Nil value for Component.Data and .Props
	NilData map[string]interface{}
	// NilClasses Nil value for Component.Classes
	NilClasses []string
	// NilAttrs Nil value for Component.Attrs
	NilAttrs map[string]string
	// NilBinds Nil value for Component.Binds
	NilBinds map[string]Bind
)

// Context -- in context component send c.Data and c.Props to method
type Context interface{}

// Method -- struct for Component methods
type Method func(Context)

// GetComponent -- function returning component child
type GetComponent func(Component) interface{}

// GetChildes -- function returning component childes
// In function parameter sends `this` component and you can get component data from this parameter
//
// Component childes can be :
// 1. String (or tag_value) like: <h1>this text is body</h1>
// 2. Another component like: <h1> <GreetingText/> </h1>
type GetChildes func(Component) []interface{}

// Bind -- struct for storage component *computed* attributes.
// It's analogue for vue `v-bind:`
// Like: `gBind:id="c.GetDataByString("iterator") + 1024"``
type Bind func(Component) string

// Component -- basic component struct
type Component struct {
	Data  map[string]interface{}
	Props map[string]interface{}

	CallBacks map[string]Method // events handlers: onClick, onHover
	Methods   map[string]Method // user functions can call from component childes

	Childes GetChildes

	Tag     string
	ID      string
	Classes []string
	Attrs   map[string]string

	Binds map[string]Bind

	// Other stuff
}

// NewComponent create new component
func NewComponent(data, props map[string]interface{}, tag string, id string, classes []string, attrs map[string]string) *Component {
	// Some stuff here, but now:
	return &Component{
		Data:  data,
		Props: props,

		Tag:     tag,
		ID:      id,
		Classes: classes,
		Attrs:   attrs,
	}
}

// AddBinds add Binds to Component
func (c *Component) AddBinds(binds map[string]Bind) *Component {
	c.Binds = binds

	return c
}

// AddChildes add childes to Component
//
// Component childes (tag value or Component) for works often requires component.Data(.Props) or .Methods
// and they (childes) can take this value from send in function Component
func (c *Component) AddChildes(childes ...GetComponent) *Component {
	c.Childes = func(this Component) []interface{} {
		var compiled []interface{}
		for _, el := range childes {
			compiled = append(compiled, el(this))
		}

		return compiled
	}

	return c
}

// CreateComponent render component.
// Not Component methods, because node can be string.
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

		_node.SetAttribute("class", strings.Join(component.Classes, " "))

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
func UpdateComponent(_parent dom.Element, new interface{}, old interface{}, index int) error {
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

		for i := 0; i < len(newChildes) || i < len(oldChildes); i++ {
			var elFromNew interface{}
			if len(newChildes) >= i {
				elFromNew = newChildes[i]
			}

			var elFromOld interface{}
			if len(oldChildes) >= i {
				elFromOld = oldChildes[i]
			}

			err = UpdateComponent(*_el, elFromNew, elFromOld, i)
			if err != nil {
				return err
			}
		}
	} else { // new or old not components => something went wrong, but we haven't errors
		return nil
	}

	return fmt.Errorf("UpdateComponent: invalid new or old types: %T, %T", new, old)
}

// GetData return data field by qury string
func (c *Component) GetData(query string) interface{} {
	// There will be callbacks, events, e.t.c.
	return c.Data[query]
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

		return newString == oldString, nil
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

		return newComponent == oldComponent, nil
	}

	return false, fmt.Errorf("changed: invalid `new` or `old`. types: %T, %T", new, old)
}
