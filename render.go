package gas

import (
	"errors"
	"fmt"
	"github.com/Sinicablyat/dom"
	"github.com/google/go-cmp/cmp"
	"log"
	"reflect"
	"strconv"
	"strings"
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

	if c.Directives.Show != nil && !c.Directives.Show(c) {
		doElHidden(_node)
	}
	_node.SetAttribute("data-i", c.UUID) // set data-i for accept element from component methods

	for attrName, attrBody := range c.Attrs {
		//if !attrIsValid(attrName, component.Tag) {continue} // check if attribute is valid for this tag
		_node.SetAttribute(attrName, attrBody)
	}

	for attrName, bindBody := range c.Binds {
		attrBody := bindBody(c)
		_node.SetAttribute(attrName, attrBody)
	}

	for handlerName, handlerBody := range c.Handlers {
		//if handlerIsValid(handlerName, component.Tag) {} // check if handler is valid for this tag

		handlerNameParsed := strings.Split(handlerName, ".")
		if len(handlerNameParsed) == 2 {
			handlerType := handlerNameParsed[0]
			handlerTarget := handlerNameParsed[1]
			switch handlerType {
			case "keyup":
				_node.AddEventListener(handlerType, func(e dom.Event) {
					if handlerTarget == strings.ToLower(e.GetValueString("key")) {
						handlerBody(c, e)
					}
				})
			case "click":
				var useTr = false
				handlerTargetInt, err := strconv.Atoi(handlerTarget)
				if err != nil {
					useTr = true
				}

				_node.AddEventListener(handlerType, func(e dom.Event) {
					buttonClick, err := parseInt(strings.ToLower(e.GetValueString("button")))
					if err != nil {
						WarnError(errors.New("invalid onClick button value"))
						return
					}

					if useTr {
						if handlerTarget == parseClickButton(buttonClick) {
							handlerBody(c, e)
						}
					} else {
						if handlerTargetInt == buttonClick {
							handlerBody(c, e)
						}
					}
				})
			}
		}
		
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


	if c.Directives.HTML.Render != nil {
		htmlDirective := c.Directives.HTML.Render(c)
		currentInner := _node.GetValueString("innerHTML")
		_node.SetInnerHTML(fmt.Sprintf("%s\n%s", currentInner, htmlDirective))
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

		// BeforeCreate hook
		if isComponent(new) && I2C(new).Hooks.BeforeCreate != nil {
			err := I2C(new).Hooks.BeforeCreate(I2C(new))
			if err != nil {
				return err
			}
		}

		_parent.AppendChild(_new)

		// Created hook
		if isComponent(new) && I2C(new).Hooks.Created != nil {
			newC := I2C(new)

			err := newC.Hooks.Created(newC)
			if err != nil {
				return err
			}
		}

		return nil
	}

	_childes := _parent.ChildNodes()
	if _childes == nil { return errors.New("_parent doesn't have childes") }

	var _el *dom.Element
	if len(_childes) > index || (len(_childes) >= index && isComponent(new)) { // component was hided if childes length <= index
		_el = _childes[index]
	}

	newC 		   := &Component{}
	newIsComponent := isComponent(new)
	if newIsComponent {newC = I2C(new)}

	if _el == nil { // here current element will exist
		return nil
	}
	
	// if component has deleted
	if new == nil {
		_parent.RemoveChild(_el)

		// Destroyed hook
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
		// BeforeUpdate hook
		if newIsComponent && I2C(new).Hooks.BeforeUpdate != nil {
			err := I2C(new).Hooks.BeforeUpdate(I2C(new))
			if err != nil {
				return err
			}
		}

		_new, err := CreateComponent(new)
		if err != nil {
			return err
		}

		_parent.ReplaceChild(_new, _el)

		// Updated hook
		if newIsComponent && I2C(new).Hooks.Updated != nil {
			err := I2C(new).Hooks.Updated(I2C(new))
			if err != nil {
				return err
			}
		}

		return nil
	}

	// check if component childes updated
	if newIsComponent {
		if newC.Directives.Model.Component != nil { // update input value
			_el.SetValue("value", newC.Directives.Model.Component.Data[newC.Directives.Model.Data])
		}

		if newC.Directives.Show != nil {
			if !newC.Directives.Show(newC) {
				doElHidden(_el)
			} else {
				doElVisible(_el)
			}
		}

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
	daE := true // cmp.Equal(new.Data, old.Data)
	wE  := cmp.Equal(new.Watchers, old.Watchers)
	mE  := true // cmp.Equal(new.Methods, old.Methods)
	coE := true // cmp.Equal(new.Computeds, old.Computeds)
	caE := cmp.Equal(new.Catchers, old.Catchers)

	hE := compareHooks(new.Hooks, old.Hooks)
	bE := compareBinds(new.renderedBinds, old.renderedBinds)

	diIfE := reflect.ValueOf(new.Directives.If).Pointer() == reflect.ValueOf(old.Directives.If).Pointer()
	diFE  := cmp.Equal(new.Directives.For, old.Directives.For)
	diME  := (new.Directives.Model.Data == old.Directives.Model.Data) && (new.Directives.Model.Component == old.Directives.Model.Component)
	diHE  := reflect.ValueOf(new.Directives.HTML.Render).Pointer() == reflect.ValueOf(old.Directives.HTML.Render).Pointer()
	diE   := diIfE && diFE && diME && diHE // Directives

	tE := new.Tag == old.Tag
	aE := cmp.Equal(new.Attrs, old.Attrs)

	return daE && wE && mE && coE && caE && hE && bE && diE && tE && aE
}

func compareHooks(new, old Hooks) bool {
	created := cmp.Equal(new.Created, old.Created)
	beforeCreate := cmp.Equal(new.BeforeCreate, old.BeforeCreate)
	destroyed := cmp.Equal(new.Destroyed, old.Destroyed)

	return created && beforeCreate && destroyed
}

func compareBinds(new, old map[string]string) bool {
	if len(new) != len(old) {
		return false
	}

	for newKey, newValue := range new {
		if newValue != old[newKey] {
			return false
		}
	}

	return true
}

// renderTree return full rendered childes tree of component
func renderTree(c *Component) []interface{} {
	var childes []interface{}
	for _, el := range c.Childes(c) {
		if isComponent(el) {
			elC := I2C(el)

			if elC.Directives.HTML.Render != nil {
				elC.Directives.HTML.Rendered = elC.Directives.HTML.Render(elC)
			}

			if elC.Binds != nil {
				if elC.renderedBinds == nil {
					elC.renderedBinds = map[string]string{}
				}

				for bindKey, bindValue := range elC.Binds { // render binds
					elC.renderedBinds[bindKey] = bindValue(elC)
				}
			}

			elC.RChildes = renderTree(elC)

			el = elC
		}

		childes = append(childes, el)
	}
	return childes
}

func doElHidden(_el *dom.Element) {
	//_el.Style().Set("visibility", "hidden")
	_el.Style().Set("display", "none")
}

func doElVisible(_el *dom.Element) {
	//_el.Style().Set("visibility", "visible")
	_el.Style().Set("display", "")
}

func parseClickButton(button int) string {
	switch button {
	case 0:
		return "left"
	case 1:
		return "middle"
	case 2:
		return "right"
	default:
		return "unknown"
	}
}

func parseInt(a string) (int, error) {
	if a == "" {
		return 0, nil
	}

	return strconv.Atoi(a)
}