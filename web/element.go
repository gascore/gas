package web

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gascore/dom"
	"github.com/gascore/gas"
)

// CreateElement render element
func CreateElement(el interface{}) (dom.Node, error) {
	switch el := el.(type) {
	case bool:
		if el {
			return createTextNode("true")
		}
		return createTextNode("false")
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return createTextNode(fmt.Sprintf("%d", el))
	case string:
		return createTextNode(el)
	case *gas.Element:
		_node, err := createHtmlElement(el)
		if err != nil {
			return nil, fmt.Errorf("cannot create component: %s", err.Error())
		}

		if el.Component != nil && el.Component.Hooks.Created != nil {
			err := el.Component.Hooks.Created()
			if err != nil {
				return nil, err
			}
		}

		if el.RefName != "" {
			p := el.ParentComponent()
			if p == nil || !p.Component.RefsAllowed {
				dom.ConsoleError("parent with allowed refs doesn't exist")
			} else {
				if p.Component.Refs == nil {
					p.Component.Refs = make(map[string]*gas.E)
				}
				p.Component.Refs[el.RefName] = el // append element to first true component refs
			}
		}

		for _, child := range el.Childes {
			_child, err := CreateElement(child)
			if err != nil {
				return nil, err
			}

			_node.AppendChild(_child)
		}

		return _node, nil
	default:
		return nil, fmt.Errorf("unsupported component type: %T", el)
	}
}

// createHtmlElement create html element without childes
func createHtmlElement(el *gas.Element) (*dom.Element, error) {
	_node := dom.NewElement(el.Tag)
	if _node == nil {
		return nil, errors.New("cannot create component")
	}

	_node.SetAttribute("data-i", el.UUID) // set data-i for accept element from component methods

	for attrName, attrBody := range el.Attrs {
		// if !attrIsValid(attrName, component.Tag) {continue} // check if attribute is valid for this tag
		_node.SetAttribute(attrName, attrBody)
	}

	for attrName, bindBody := range el.Binds {
		_node.SetAttribute(attrName, bindBody())
	}

	for handlerName, handlerBody := range el.Handlers {
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
						handlerBody(ToUniteObject(e))
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
					e.PreventDefault() // Because i don't want trigger <a>

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
							handlerBody(ToUniteObject(e))
						}

					} else {
						if handlerTargetInt == buttonClick {
							handlerBody(ToUniteObject(e))
						}
					}
				})

				continue
			}
		} else {
			_node.AddEventListener(handlerName, func(e dom.Event) {
				handlerBody(ToUniteObject(e))
			})
		}
	}

	if len(el.Watcher) != 0 && (el.Tag == "input" || el.Tag == "textarea" || el.Tag == "select") { // model allowed only for input tags
		p := el.ParentComponent()
		c := p.Component
		if c.Watchers == nil || len(c.Watchers) == 0 {
			return nil, fmt.Errorf("invalid watcher for \"%s\" in component \"%s\"", el.Watcher, el.UUID)
		}

		watcher, ok := c.Watchers[el.Watcher]
		if !ok {
			return nil, fmt.Errorf("watcher \"%s\" is undefined in component \"%s\"", el.Watcher, p.UUID)
		}

		startVal, err := watcher(nil, nil)
		if err != nil {
			return nil, err
		}

		_node.SetValue(startVal)

		_node.AddEventListener("input", func(e dom.Event) {
			_target := e.Target()
			inputValue := _target.Value()

			var (
				// inputType       string
				inputValueTyped interface{}
				err             error
			)

			switch el.Attrs["type"] {
			case "range", "number":
				// inputType = "int"
				inputValueTyped, err = strconv.Atoi(inputValue)
				if err != nil {
					warnError(err)
					return
				}
			case "checkbox":
				// inputType = "bool"
				inputValueTyped = _target.JSValue().Get("checked").Bool()
			default:
				// inputType = "string"
				inputValueTyped = inputValue
			}

			newVal, err := watcher(inputValueTyped, ToUniteObject(e))
			if err != nil {
				warnError(err)
			} else {
				_node.SetValue(newVal)
				go c.Update()
			}
		})
	}

	if el.HTML.Render != nil {
		htmlDirective := el.HTML.Render()
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

func (w BackEnd) EditWatcherValue(el interface{}, newVal string) {
	el.(*dom.Element).SetValue(newVal)
}
