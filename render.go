package gas

import(
	"fmt"
	"errors"
	"github.com/Sinicablyat/dom"
)

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

			_node.AppendChild(_child)
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

		_parent.AppendChild(_new)

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

		_parent.ReplaceChild(_new, _el)

		return nil
	}

	// check if component childes updated
	if newIsComponent {
		_new, err := CreateElement(newC)
		if err != nil {
			return err
		}

		for len(_el.ChildNodes()) > 0 { // transfer _old childes to _new
			_new.AppendChild(_el.ChildNodes()[0])
		}

		_parent.ReplaceChild(_new, _el)
		
		// update childes
		newChildes := newC.Childes(newC) // new.Childes(new)
		oldChildes := I2C(old).Childes(I2C(old)) // old.Childes(old)

		err = UpdateComponentChildes(_new, newChildes, oldChildes)
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