package web

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gascore/dom"
	"github.com/gascore/gas"
)

// CreateComponent render component. Returns _el, err
func CreateComponent(component interface{}) (dom.Node, error) {
	switch component := component.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return createTextNode(fmt.Sprintf("%d", component))
	case string:
		return createTextNode(component)
	case *gas.Component:
		_node, err := CreateElement(component)
		if err != nil {
			return nil, errors.New("cannot create component")
		}

		if component.Hooks.Created != nil {
			err := component.Hooks.Created(component)
			if err != nil {
				return nil, err
			}
		}

		component.UpdateHTMLDirective()

		if component.Ref != "" {
			p := component.ParentWithAllowedRefs()
			dom.ConsoleLog(p == nil)
			if p == nil {
				dom.ConsoleError("parent with allowed refs doesn't exist")
			} else {
				if p.Refs == nil {
					p.Refs = make(map[string]*gas.C)
				}
				p.Refs[component.Ref] = component // append element to first true component refs
			}
		}

		for _, el := range component.RChildes {
			_child, err := CreateComponent(el)
			if err != nil {
				return nil, err
			}

			_node.AppendChild(_child)
		}

		return _node, nil
	default:
		return nil, fmt.Errorf("invalid component type: %T", component)
	}
}

// CreateElement create html element without childes
func CreateElement(c *gas.Component) (*dom.Element, error) {
	_node := dom.NewElement(c.Tag)
	if _node == nil {
		return nil, errors.New("cannot create component")
	}

	updateVisible(c, _node)

	_node.SetAttribute("data-i", c.UUID) // set data-i for accept element from component methods

	for attrName, attrBody := range c.Attrs {
		// if !attrIsValid(attrName, component.Tag) {continue} // check if attribute is valid for this tag
		_node.SetAttribute(attrName, attrBody)
	}

	for attrName, bindBody := range c.Binds {
		_node.SetAttribute(attrName, bindBody())
	}

	for handlerName, handlerBody := range c.Handlers {
		// if handlerIsValid(handlerName, component.Tag) {} // check if handler is valid for this tag

		handlerNameParsed := strings.Split(handlerName, ".")
		if len(handlerNameParsed) == 2 {
			handlerType := handlerNameParsed[0]
			handlerTarget := handlerNameParsed[1]
			switch handlerType {
			case "keyup":
				_node.AddEventListener("keyup", func(e dom.Event) {
					if handlerTarget == strings.ToLower(e.Key()) ||
						handlerTarget == strings.ToLower(e.KeyCode()) {
						handlerBody(c, ToUniteObject(e))
					}
				})
				continue
			case "click":
				var useTr = false
				handlerTargetInt, err := strconv.Atoi(handlerTarget)
				if err != nil {
					useTr = true
				}

				_node.AddEventListener("click", func(e dom.Event) {
					buttonClick, err := parseInt(strings.ToLower(e.ButtonAttr()))
					if err != nil {
						return
					}

					if useTr {
						var parsedButtonClick string
						switch buttonClick {
						case 0:
							parsedButtonClick = "left"
						case 1:
							parsedButtonClick = "middle"
						case 2:
							parsedButtonClick = "right"
						default:
							parsedButtonClick = "unknown"
						}

						if handlerTarget == parsedButtonClick {
							handlerBody(c, ToUniteObject(e))
						}

					} else {
						if handlerTargetInt == buttonClick {
							handlerBody(c, ToUniteObject(e))
						}
					}
				})

				continue
			}
		} else {
			_node.AddEventListener(handlerName, func(e dom.Event) {
				handlerBody(c, ToUniteObject(e))
			})
		}
	}

	if len(c.Model.Data) != 0 && (c.Tag == "input" || c.Tag == "textarea" || c.Tag == "select") { // model allowed only for input tags
		this := c.Model.Component
		dataS := c.Model.Data
		deep := c.Model.Deep

		if len(deep) == 0 {
			_node.SetValue(this.Get(dataS))
		} else {
			dataValue := this.Get(dataS)
			field, err := getField(reflect.ValueOf(dataValue), deep)
			if err != nil {
				return nil, err
			}

			_node.SetValue(field.Interface())
		}

		_node.AddEventListener("input", func(e dom.Event) {
			_target := e.Target()
			inputValue := _target.Value()

			var (
				inputType       string
				inputValueTyped interface{}
				err             error
			)

			switch c.Attrs["type"] {
			case "range", "number":
				inputType = "int"
				inputValueTyped, err = strconv.Atoi(inputValue)
				if err != nil {
					warnError(err)
					return
				}
			case "checkbox":
				inputType = "bool"
				inputValueTyped = _target.JSValue().Get("checked").Bool()
			default:
				inputType = "string"
				inputValueTyped = inputValue
			}

			dataValue := this.Get(dataS)

			if len(deep) != 0 {
				err := (func() error {
					field, err := getField(reflect.ValueOf(dataValue), deep)
					if err != nil {
						return err
					}

					switch inputValueTyped.(type) {
					case int:
						field.SetInt(int64(inputValueTyped.(int)))
					case string:
						field.SetString(inputValue)
					case bool:
						field.SetBool(inputValueTyped.(bool))
					}

					return nil
				})()
				if err != nil {
					warnError(err)
					return
				}
				inputValueTyped = dataValue
			} else {
				if inputType != reflect.TypeOf(dataValue).String() {
					warnError(errors.New("input type != data type"))
					return
				}
			}

			c.RC.Add(singleNode(&gas.RenderNode{
				Type:     gas.DataType,
				Priority: gas.InputPriority,
				New:      this,
				Data: map[string]interface{}{
					dataS: inputValueTyped,
				},
			}))
		})
	}

	if c.HTML.Render != nil {
		htmlDirective := c.HTML.Render(c)
		_node.SetInnerHTML(fmt.Sprintf("%s\n%s", _node.InnerHTML(), htmlDirective))
	}

	return _node, nil
}

// createTextNode create TextNode by string
func createTextNode(node string) (dom.Node, error) {
	_node := dom.Doc.CreateTextNode(node)
	if _node == nil {
		return nil, errors.New("cannot create textNode")
	}

	return _node, nil
}

func singleNode(node *gas.RenderNode) []*gas.RenderNode {
	return []*gas.RenderNode{node}
}
