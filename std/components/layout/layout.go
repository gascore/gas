package layout

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gascore/dom"
	"github.com/gascore/dom/js"
	"github.com/gascore/gas"
	sjs "syscall/js"
)

// Config layout config structure
type Config struct {
	DragInterval int

	LayoutClass string
	GutterClass string
	GutterSize  int

	Sizes []Size

	Type bool // true - horizontal, false - vertical

	OnStart Event
	OnStop  Event // runs before Element recreating
	OnMove  MoveEvent

	byGuttersOffset float64
	allGuttersSize  int

	typeString     string
	orientation    string
	orientationB   string
	subOrientation string
	clientAxis     string
	positionEnd    string
}

type Event func(first, second Element, _gutter *dom.Element) (stopIt bool, err error)
type MoveEvent func(first, second Element, _gutter *dom.Element, offset float64) (stopIt bool, err error)

// Size layout item size info
type Size struct {
	Min     float64
	Max     float64
	Start   float64
	current float64
}

type Element struct {
	C     *gas.Component
	Index int
}

// Layout return resizable layout
func Layout(config *Config, e gas.External) *gas.Component {
	e.Body = gas.RemoveStrings(e.Body)

	if len(e.Body) != len(config.Sizes) {
		dom.ConsoleError("not enough Element sizes")
		return nil
	}

	if config.DragInterval == 0 {
		config.DragInterval = 1
	}

	if config.GutterSize == 0 {
		config.GutterSize = 2
	}

	if config.Type {
		config.typeString = "horizontal"
		config.orientation = "width"
		config.orientationB = "Width"
		config.subOrientation = "height"
		config.clientAxis = "clientX"
		config.positionEnd = "right"
	} else {
		config.typeString = "vertical"
		config.orientation = "height"
		config.orientationB = "Height"
		config.subOrientation = "width"
		config.clientAxis = "clientY"
		config.positionEnd = "bottom"
	}

	if config.LayoutClass == "" {
		config.LayoutClass = "layout"
	}

	if config.GutterClass == "" {
		config.GutterClass = "gutter"
	}

	var sizesSum float64 // check sizes sum == 100 && make Size.Start valid
	for _, size := range config.Sizes {
		if size.Start > size.Max {
			size.Start = size.Max
		} else if size.Min > size.Start {
			size.Start = size.Min
		}
		sizesSum += size.Start
	}

	if sizesSum != 100 {
		// because i don't want to create an implicit state change
		dom.ConsoleError("invalid sizes: size.Start sum != 100")
		return nil
	}

	var elements []Element
	var sizes []Size
	for i, child := range e.Body {
		if !gas.IsComponent(child) {
			dom.ConsoleError(fmt.Sprintf("invalid child in layout - child is not component, want: '*gas.Component' got: '%T'", child))
			return nil
		}

		childC := gas.I2C(child)

		if childC.Attrs == nil {
			childC.Attrs = make(map[string]string)
		}

		size := config.Sizes[i]
		size.current = size.Start

		sizes = append(sizes, size)
		elements = append(elements, Element{C: childC, Index: i})
	}

	config.allGuttersSize = (len(elements) - 1) * config.GutterSize
	config.byGuttersOffset = float64(config.allGuttersSize) / float64(len(elements))

	return gas.NC(
		&gas.C{
			Data: map[string]interface{}{
				"sizes": sizes,
			},
			Attrs: map[string]string{
				"class": fmt.Sprintf("%s %s-%s", config.LayoutClass, config.LayoutClass, config.typeString),
			},
		},
		func(this *gas.Component) []interface{} {
			var childes []interface{}

			aSizes := this.Get("sizes").([]Size)

			for i, child := range elements {
				childes = append(childes, gas.NE(
					&gas.C{
						Attrs: map[string]string{
							"class":  config.LayoutClass + "-item",
							"style":  fmt.Sprintf("%s: calc(%f%s - %fpx); %s: 100%s;", config.orientation, aSizes[child.Index].current, "%", config.byGuttersOffset, config.subOrientation, "%"),
							"data-i": fmt.Sprintf("%d", i),
						},
					},
					child.C,
				))
				if i != len(elements)-1 {
					childes = append(childes, gutter(this, config, child, elements[i+1]))
				}
			}

			return childes
		},
	)
}

func gutter(pThis *gas.C, config *Config, first, second Element) *gas.Component {
	var cursorType string
	if config.Type {
		cursorType = "ew-resize"
	} else {
		cursorType = "row-resize"
	}

	return gas.NE(&gas.C{
		Data: map[string]interface{}{
			"dragOffset": float64(0),
			"dragging":   false,

			"startEvent": js.Func{},
			"moveEvent":  js.Func{},
			"stopEvent":  js.Func{},
		},
		Attrs: map[string]string{
			"class": fmt.Sprintf("%s %s-%s", config.GutterClass, config.GutterClass, config.typeString),
			"style": fmt.Sprintf("cursor: %s; %s: %dpx", cursorType, config.orientation, config.GutterSize),
		},
		Hooks: gas.Hooks{
			Mounted: func(this *gas.Component) error {
				var parentSize int

				_el := this.Element().(*dom.Element)
				computedStyles := sjs.Global().Call("getComputedStyle", _el.JSValue())

				if config.Type {
					parentSize = _el.ParentElement().ClientHeight() - parseP(computedStyles.Get("paddingTop")) - parseP(computedStyles.Get("paddingBottom"))
				} else {
					parentSize = _el.ParentElement().ClientWidth() - parseP(computedStyles.Get("paddingLeft")) - parseP(computedStyles.Get("paddingRight"))
				}

				moveEvent := event(func(event dom.Event) error {
					if !this.Get("dragging").(bool) {
						return nil
					}

					event.PreventDefault()

					var start float64
					if config.Type {
						start = _el.GetBoundingClientRectRaw().Get("left").Float()
					} else {
						start = _el.GetBoundingClientRectRaw().Get("top").Float()
					}

					offset := getMousePosition(config.clientAxis, event) - start + float64(config.GutterSize)
					if offset == 0 {
						return nil
					}

					sizes := pThis.Get("sizes").([]Size)

					var newFirst, newSecond float64
					if offset < 0 {
						newFirst, newSecond = getSizes(-offset, _el.ParentElement(), sizes[first.Index], sizes[second.Index], config)
					} else {
						newSecond, newFirst = getSizes(offset, _el.ParentElement(), sizes[second.Index], sizes[first.Index], config)
					}

					sizes[first.Index].current = newFirst
					sizes[second.Index].current = newSecond

					pThis.SetValue("sizes", sizes)

					if config.OnMove != nil {
						stopIt, err := config.OnMove(first, second, _el, offset)
						if err != nil {
							return err
						}

						if stopIt {
							return nil
						}
					}

					return nil
				})

				stopEvent := event(func(event dom.Event) error {
					if !this.Get("dragging").(bool) {
						return nil
					}

					_first := first.C.Element().(*dom.Element)
					_second := second.C.Element().(*dom.Element)

					if config.OnStop != nil {
						stopIt, err := config.OnStop(first, second, _el)
						if err != nil {
							return err
						}

						if stopIt {
							return nil
						}
					}

					removeEvent(_el, "touchend", this.Get("stopEvent").(js.Func))
					removeEvent(_el, "touchcancel", this.Get("stopEvent").(js.Func))
					removeEvent(_el, "touchmove", this.Get("moveEvent").(js.Func))
					removeEvent(dom.Doc, "mouseup", this.Get("stopEvent").(js.Func))
					removeEvent(dom.Doc, "mousemove", this.Get("moveEvent").(js.Func))

					for _, _x := range []*dom.Element{_first, _second} {
						_x.Style().Set("userSelect", "")
						_x.Style().Set("webkitUserSelect", "")
						_x.Style().Set("MozUserSelect", "")
					}

					this.Set(map[string]interface{}{
						"dragging": false,
					})

					return nil
				})

				startEvent := event(func(event dom.Event) error {
					if this.Get("dragging").(bool) {
						return nil
					}

					_el := event.Target()
					_first := first.C.Element().(*dom.Element)
					_second := second.C.Element().(*dom.Element)

					if config.OnStart != nil {
						stopIt, err := config.OnStart(first, second, _el)
						if err != nil {
							return err
						}

						if stopIt {
							return nil
						}
					}

					addEvent(_el, "touchend", this.Get("stopEvent").(js.Func))
					addEvent(_el, "touchcancel", this.Get("stopEvent").(js.Func))
					addEvent(_el, "touchmove", this.Get("moveEvent").(js.Func))
					addEvent(dom.Doc, "mouseup", this.Get("stopEvent").(js.Func))
					addEvent(dom.Doc, "mousemove", this.Get("moveEvent").(js.Func))

					for _, _x := range []*dom.Element{_first, _second} {
						_x.Style().Set("userSelect", "none")
						_x.Style().Set("webkitUserSelect", "none")
						_x.Style().Set("MozUserSelect", "none")
					}

					_el.ClassList().Add(config.GutterClass + "-focus")

					this.Set(map[string]interface{}{
						"dragOffset": getMousePosition(config.clientAxis, event) - _first.GetBoundingClientRectRaw().Get(config.positionEnd).Float(),
						"dragging":   true,
					})

					return nil
				})

				_el.Style().Set(config.subOrientation, fmt.Sprintf("%dpx", parentSize))

				addEvent(_el, "mousedown", startEvent)
				addEvent(_el, "touchstart", startEvent)
				addEvent(_el, "touchend", stopEvent)
				addEvent(_el, "touchcancel", stopEvent)
				addEvent(_el, "touchmove", moveEvent)

				addEvent(_el, "mouseup", stopEvent)
				addEvent(_el, "mousemove", moveEvent)

				this.SetImm(map[string]interface{}{
					"startEvent": startEvent,
					"moveEvent":  moveEvent,
					"stopEvent":  stopEvent,
				})

				return nil
			},
			BeforeDestroy: func(this *gas.C) error {
				_el := this.Element().(*dom.Element)
				if _el == nil {
					return nil
				}

				removeEvent(_el, "mousedown", this.Get("startEvent").(js.Func))
				removeEvent(_el, "touchstart", this.Get("startEvent").(js.Func))

				return nil
			},
		},
	})
}

func event(f func(event dom.Event) error) js.Func {
	return js.NewEventCallback(func(v js.Value) {
		err := f(dom.ConvertEvent(v))
		if err != nil {
			dom.ConsoleError(err.Error())
			return
		}
	})
}

func addEvent(e dom.Node, typ string, h js.Func) {
	e.JSValue().Call("addEventListener", typ, h)
}

func removeEvent(e dom.Node, typ string, h js.Func) {
	e.JSValue().Call("removeEventListener", typ, h)
}

func getSizes(offset float64, _parent *dom.Element, first, second Size, config *Config) (float64, float64) {
	offsetOri := "offset" + config.orientationB
	layoutSize := _parent.JSValue().Get(offsetOri).Float() - float64(config.allGuttersSize)
	theirSizeP := first.current + second.current
	offsetP := (offset * 100) / layoutSize

	if first.current-offsetP < first.Min {
		makeSecond := theirSizeP - first.Min
		if makeSecond > second.Max {
			return theirSizeP - second.Max, second.Max
		}

		return first.Min, makeSecond
	}

	if second.current+offsetP > second.Max {
		if second.Max > theirSizeP || second.Max > theirSizeP-first.Min {
			return first.Min, theirSizeP - first.Min
		}

		return theirSizeP - second.Max, second.Max
	}

	return first.current - offsetP, theirSizeP - (first.current - offsetP)
}

func getMousePosition(clientAxis string, event dom.Event) float64 {
	if notJsNull(event.JSValue().Get("touches")) {
		return event.JSValue().Get("touches").Get("0").Get(clientAxis).Float()
	}

	return event.JSValue().Get(clientAxis).Float()
}

func notJsNull(e sjs.Value) bool {
	return e.Type() != sjs.TypeUndefined && e.Type() != sjs.TypeNull
}

func parseP(a sjs.Value) int {
	b, _ := strconv.Atoi(strings.TrimSuffix(a.String(), "px"))
	return b
}
