package wasm

import (
	"errors"
	"fmt"
	"github.com/Sinicablyat/dom"
	"github.com/Sinicablyat/gas/core"
	"strconv"
	"strings"
)

// nil backend. Need for call out-methods
var be BackEnd

// CreateComponent render component. Returns _el, err
func (w BackEnd) CreateComponent(node interface{}) (*dom.Element, error) {
	switch node.(type) {
	case string:
		nodeS := node.(string)
		_node := dom.NewElement("span") // dennwc/dom doesn't support textNode

		if _node == nil {
			return nil, errors.New("cannot create textNode")
		}

		_node.SetTextContent(nodeS)

		return _node, nil
	case *core.Component:
		component := core.I2C(node)

		_node, err := CreateElement(component)
		if err != nil {
			return nil, errors.New("cannot create component")
		}

		for _, el := range component.Childes(component) {
			_child, err := w.CreateComponent(el)
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
func CreateElement(c *core.Component) (*dom.Element, error) {
	_node := dom.NewElement(c.Tag)
	if _node == nil {
		return nil, errors.New("cannot create component")
	}

	setShow(c, _node)

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
						//log.Error("invalid onClick button value")
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

	if len(c.Directives.Model.Data) != 0 && (c.Tag == "input" || c.Tag == "textarea" || c.Tag == "select") { // model allowed only for <input>
		this  := c.Directives.Model.Component
		dataS := c.Directives.Model.Data

		_node.SetValue("value", this.GetData(dataS))
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
					//WarnError(err)
					return
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

			dataValue := this.GetData(dataS)
			if inputType != fmt.Sprintf("%T", dataValue) {
				//WarnError(errors.New("input type != data type"))
				return
			}

			err = this.SetData(dataS, inputValueTyped)
			if err != nil {
				//WarnError(err)
				return
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


// UpdateComponent trying to Update component
func UpdateComponent(_parent *dom.Element, new interface{}, old interface{}, index int) error {
	// if component has created
	if old == nil {
		_new, err := be.CreateComponent(new)
		if err != nil {
			return err
		}

		// BeforeCreate hook
		if core.IsComponent(new) && core.I2C(new).Hooks.BeforeCreate != nil {
			err := core.I2C(new).Hooks.BeforeCreate(core.I2C(new))
			if err != nil {
				return err
			}
		}

		_parent.AppendChild(_new)

		// Created hook
		if core.IsComponent(new) && core.I2C(new).Hooks.Created != nil {
			newC := core.I2C(new)

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
	if len(_childes) > index { // component was hided if childes length <= index
		_el = _childes[index]
	}

	newC 		   := &core.Component{}
	newIsComponent := core.IsComponent(new)
	if newIsComponent {newC = core.I2C(new)}

	if _el == nil { // here current element will exist
		return nil
	}
	
	// if component has deleted
	if new == nil {
		_parent.RemoveChild(_el)

		// Destroyed hook
		if core.IsComponent(old) && core.I2C(old).Hooks.Destroyed != nil {
			err := core.I2C(old).Hooks.Destroyed(core.I2C(old))
			if err != nil {
				return err
			}
		}

		return nil
	}

	// if component has Changed
	isChanged, err := core.Changed(new, old)
	if err != nil {
		return err
	}
	if isChanged {
		// BeforeUpdate hook
		if newIsComponent && core.I2C(new).Hooks.BeforeUpdate != nil {
			newC := core.I2C(new)
			err := newC.Hooks.BeforeUpdate(newC)
			if err != nil {
				return err
			}
		}

		_new, err := be.CreateComponent(new)
		if err != nil {
			return err
		}

		_parent.ReplaceChild(_new, _el)

		// Updated hook
		if newIsComponent && newC.Hooks.Updated != nil {
			err := newC.Hooks.Updated(newC)
			if err != nil {
				return err
			}
		}

		return nil
	}

	// check if component childes updated
	if newIsComponent {
		if newC.Directives.Model.Component != nil { // Update input value
			_el.SetValue("value", newC.Directives.Model.Component.Data[newC.Directives.Model.Data])
		}

		setShow(newC, _el)

		err = DeepUpdateComponentChildes(_el, newC.RChildes, core.I2C(old).RChildes)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateComponentChildes start point for update component childes from core
func (w BackEnd) UpdateComponentChildes(c *core.Component, newChildes, oldChildes []interface{}) error {
	_el := c.GetElement().(*dom.Element)
	return DeepUpdateComponentChildes(_el, newChildes, oldChildes)
}

// DeepUpdateComponentChildes Update component childes by new and old childes
func DeepUpdateComponentChildes(_el *dom.Element, newChildes, oldChildes []interface{}) error {
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


// RenderTree return full rendered childes tree of component
func (w BackEnd) RenderTree(c *core.Component) []interface{} {
	var childes []interface{}
	for _, el := range c.Childes(c) {
		if core.IsComponent(el) {
			elC := core.I2C(el)

			if elC.Directives.HTML.Render != nil {
				elC.Directives.HTML.Rendered = elC.Directives.HTML.Render(elC)
			}

			if elC.Binds != nil {
				if elC.RenderedBinds == nil {
					elC.RenderedBinds = map[string]string{}
				}

				for bindKey, bindValue := range elC.Binds { // render binds
					elC.RenderedBinds[bindKey] = bindValue(elC)
				}
			}

			elC.RChildes = w.RenderTree(elC)

			el = elC
		}

		childes = append(childes, el)
	}
	return childes
}

func (w BackEnd) ReCreate(c *core.Component) error {
	_updatedC, err := w.CreateComponent(c)
	if err != nil {
		return err
	}

	c.ParentC.GetElement().(*dom.Element).ReplaceChild(_updatedC, c.GetElement().(*dom.Element))

	return nil
}


// GetElement get dom.Element by component
func (w BackEnd) GetElement(c *core.Component) interface{} {
	return dom.Doc.QuerySelector(fmt.Sprintf("[data-i='%s']", c.UUID)) // select element by data-i attribute
}

// GetElement get dom.Element by component
func (w BackEnd) GetGasEl(g *core.Gas) interface{} {
	return dom.Doc.GetElementById(g.StartPoint)
}


func setShow(c *core.Component, _el *dom.Element) {
	if c.Directives.Show != nil {
		if !c.Directives.Show(c) {
			doElHidden(_el)
		} else {
			doElVisible(_el)
		}
	}
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