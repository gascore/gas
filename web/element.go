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
	case string:
		return createTextNode(el)
	case fmt.Stringer:
		return createTextNode(el.String())
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return createTextNode(fmt.Sprintf("%v", el))
	case bool:
		if el {
			return createTextNode("true")
		}
		return createTextNode("false")
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

func newNodeEl(tag string) *dom.Element {
	switch tag {
	case "animate", "animatemotion", "animatetransform", "circle", "clippath", "color-profile", "defs", "desc", "discard", "ellipse", "feblend", "fecolormatrix", "fecomponenttransfer", "fecomposite", "feconvolvematrix", "fediffuselighting", "fedisplacementmap", "fedistantlight", "fedropshadow", "feflood", "fefunca", "fefuncb", "fefuncg", "fefuncr", "fegaussianblur", "feimage", "femerge", "femergenode", "femorphology", "feoffset", "fepointlight", "fespecularlighting", "fespotlight", "fetile", "feturbulence", "filter", "foreignobject", "g", "hatch", "hatchpath", "image", "line", "lineargradient", "marker", "mask", "mesh", "meshgradient", "meshpatch", "meshrow", "metadata", "mpath", "path", "pattern", "polygon", "polyline", "radialgradient", "rect", "set", "solidcolor", "stop", "svg", "switch", "symbol", "text", "textpath", "title", "tspan", "unknown", "use", "view", "svg-style", "svg-script", "svg-a":
		return dom.NewElementNS("http://www.w3.org/2000/svg", strings.TrimPrefix(tag, "svg-"))
	default:
		return dom.NewElement(tag)
	}
}

// createHtmlElement create html element without children
func createHtmlElement(el *gas.Element) (*dom.Element, error) {
	_node := newNodeEl(el.Tag)
	if _node == nil {
		return nil, errors.New("cannot create component")
	}

	_node.SetAttribute("data-i", el.UUID) // set data-i for accept element from component methods

	setAttributes(_node, el.RAttrs)

	for handlerName, handlerBodyG := range el.Handlers {
		handlerBody := handlerBodyG
		handlerNameParsed := strings.Split(handlerName, ".")
		if len(handlerNameParsed) == 2 {
			handlerType := handlerNameParsed[0]
			handlerTarget := strings.ToLower(handlerNameParsed[1])

			var handlerTargetIsInt bool
			handlerTargetInt, err := strconv.Atoi(handlerTarget)
			handlerTargetIsInt = err == nil

			switch handlerType {
			case "keyup":
				_node.AddEventListener("keyup", func(e dom.Event) {
					keyCode, _ := strconv.Atoi(e.KeyCode())
					if handlerTarget == strings.ToLower(e.Key()) || (handlerTargetIsInt && handlerTargetInt == keyCode) {
						handlerBody(ToGasEvent(e, false))
					}
				})
				continue
			case "click":
				_node.AddEventListener("click", func(e dom.Event) {
					e.PreventDefault() // Because i don't want trigger <a>

					buttonClick, err := parseInt(strings.ToLower(e.ButtonAttr()))
					if err != nil {
						return
					}

					var correctKey bool
					if handlerTargetIsInt {
						correctKey = handlerTargetInt == buttonClick
					} else {
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

						correctKey = handlerTarget == parsedButtonClick
					}

					if correctKey {
						handlerBody(ToGasEvent(e, false))
					}
				})

				continue
			}
		}

		isCheckbox := isInputCheckbox(el.Tag, el.RAttrs)
		_node.AddEventListener(handlerName, func(e dom.Event) {
			handlerBody(ToGasEvent(e, isCheckbox))
		})
	}

	if el.HTML != nil {
		_node.SetInnerHTML(_node.InnerHTML() + "\n" + el.RHTML)
	}

	return _node, nil
}

func setAttributes(_el *dom.Element, attrs gas.Map) {
	for attrKey, attrVal := range attrs {
		if attrKey == "checked" {
			if attrVal == "false" {
				_el.JSValue().Set("checked", false)
				_el.RemoveAttribute("checked")
				continue
			}

			_el.JSValue().Set("checked", true)
		}
		
		_el.SetAttribute(attrKey, attrVal)
		
		if attrKey == "value" {
			_el.SetValue(attrVal)
		}
	}
}

func isInputCheckbox(tag string, attrs gas.Map) bool {
	if tag != "input" {
		return false
	}

	for attrKey, attrVal := range attrs {
		if attrKey == "type" && attrVal == "checkbox" {
			return true
		}
	}

	return false
}

// createTextNode create TextNode by string
func createTextNode(node string) (dom.Node, error) {
	_node := dom.Doc.CreateTextNode(node)
	if _node == nil {
		return nil, errors.New("cannot create textNode")
	}

	return _node, nil
}
