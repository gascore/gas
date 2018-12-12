package gas

import (
	"errors"
	"fmt"
	"github.com/Sinicablyat/dom"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"strconv"
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
			handlerBody(c, e)
		})
	}

	if len(c.Directives.Model.Data) != 0 && c.Tag == "input" { // model allowed only for <input>
		_node.AddEventListener("input", func(e dom.Event) {
			_target := e.Target()
			inputValue := _target.GetValue("value").String()

			var (
				inputType string
				inputValueTyped interface{}
				err error
			)
			switch c.Attrs["type"] {
			case "range", "number":
				inputType = "int"
				inputValueTyped, err = strconv.Atoi(inputValue)
				if err != nil {
					WarnError(err)
				}

				break
			case "checkbox":
				inputType = "bool"
				inputValueTyped = _target.GetValue("checked").Bool()
				break
			default:
				inputType = "string"
				inputValueTyped = inputValue
				break
			}

			this := c.Directives.Model.Component

			dataValue := this.GetData(c.Directives.Model.Data)
			if inputType != fmt.Sprintf("%T", dataValue) {
				WarnError(errors.New("input type != data type"))
			}

			err = this.SetData(c.Directives.Model.Data, inputValueTyped)
			if err != nil {
				WarnError(err)
			}
		})
	}

	htmlDirective := c.Directives.HTML.Render(c)
	if len(htmlDirective) != 0 {
		currentInner := _node.GetValueString("innerHTML")
		_node.SetInnerHTML(fmt.Sprintf("%s\n%s", currentInner, htmlDirective))
	}

	if c.Hooks.BeforeCreate != nil {
		err := c.Hooks.BeforeCreate(c)
		if err != nil {
			return nil, err
		}
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

		if isComponent(new) && I2C(new).Hooks.Created != nil {
			newC := I2C(new)

			err := newC.eventInUpdater(func() error {
				return newC.Hooks.Created(newC)
			})
			if err != nil {
				return err
			}
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
	if newIsComponent {newC = I2C(new)}

	// if component has deleted
	if new == nil {
		_parent.RemoveChild(_el)

		if isComponent(old) && I2C(old).Hooks.Destroyed != nil {
			err := I2C(old).Hooks.Destroyed(I2C(old))
			if err != nil {
				return err
			}
		}

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
		err = UpdateComponentChildes(_el, newC.RChildes, I2C(old).RChildes)
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
		if len(newChildes) > i {
			elFromNew = newChildes[i]
		}

		var elFromOld interface{}
		if len(oldChildes) > i {
			elFromOld = oldChildes[i]
		}

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
		newC := I2C(new)
		oldC := I2C(old)

		if newC.Directives.HTML.Rendered != oldC.Directives.HTML.Rendered {
			return true, nil
		}

		return !isComponentsEquals(newC, oldC), nil // thank you god for the go-cmp
	}

	return false, fmt.Errorf("changed: invalid `new` or `old`. types: %T, %T", new, old)
}

func isComponentsEquals(new, old *Component) bool {
	// sometimes i'm sad that i chose strict-type pl
	daE := cmp.Equal(new.Data, old.Data)
	wE := cmp.Equal(new.Watchers, old.Watchers)
	//mE := cmp.Equal(new.Methods, old.Methods)
	mE := true
	//coE := cmp.Equal(new.Computeds, old.Computeds)
	coE := true
	caE := cmp.Equal(new.Catchers, old.Catchers)
	hE := cmp.Equal(new.Handlers, old.Handlers)
	bE := cmp.Equal(new.Binds, old.Binds)

	diIfE := reflect.ValueOf(new.Directives.If).Pointer() == reflect.ValueOf(old.Directives.If).Pointer()
	diFE := cmp.Equal(new.Directives.For, old.Directives.For)
	diME := cmp.Equal(new.Directives.Model, old.Directives.Model)
	diHE := reflect.ValueOf(new.Directives.HTML.Render).Pointer() == reflect.ValueOf(old.Directives.HTML.Render).Pointer()
	diE := diIfE && diFE && diME && diHE // Directives

	tE := new.Tag == old.Tag
	aE := cmp.Equal(new.Attrs, old.Attrs)

	return daE && wE && mE && coE && caE && hE && bE && diE && tE && aE
}

// renderTree return full rendered childes tree of component
func renderTree(c *Component) []interface{} {
	var childes []interface{}
	for _, el := range c.Childes(c) {
		if isComponent(el) {
			elC := I2C(el)

			elC.Directives.HTML.Rendered = elC.Directives.HTML.Render(elC)
			elC.RChildes = renderTree(elC)

			el = elC
		}

		childes = append(childes, el)
	}
	return childes
}