package dndfree

import (
	"errors"
	"fmt"

	"github.com/frankenbeanies/uuid4"

	web "github.com/gascore/gas/web"

	"github.com/gascore/dom"
	"github.com/gascore/dom/js"
	"github.com/gascore/gas"
)

// Config dnd-free config structure
type Config struct {
	Tag string

	XDisabled bool // only Y
	YDisabled bool // only X

	ifXDisabled int // 0 if disabled, 1 if enabled
	ifYDisabled int // 0 if disabled, 1 if enabled

	Boundary string // boundary element class
	Handle   string // handle element class

	Class string

	OnMove  func(event gas.Object, x, y int) error
	OnStart func(event gas.Object, x, y int) (block bool, err error)
	OnEnd   func(event gas.Object, x, y int) (reset bool, err error)
}

// DNDFree free draggable component
func DNDFree(config Config, e gas.External) *gas.C {
	if config.Tag == "" {
		config.Tag = "div"
	}

	if config.Class == "" {
		config.Class = "dnd-free"
	}

	if config.XDisabled && config.YDisabled {
		dom.ConsoleError("x and y are disabled: element is static")
		return nil
	}

	if config.XDisabled {
		config.ifXDisabled = 0
	} else {
		config.ifXDisabled = 1
	}

	if config.YDisabled {
		config.ifYDisabled = 0
	} else {
		config.ifYDisabled = 1
	}

	childUUID := uuid4.New().String()

	return gas.NC(
		&gas.C{
			Tag: config.Tag,
			Data: map[string]interface{}{
				"initialX": 0,
				"initialY": 0,

				"offsetX": 0,
				"offsetY": 0,

				"cursorOffsetLeft":   0,
				"cursorOffsetTop":    0,
				"cursorOffsetRight":  0,
				"cursorOffsetBottom": 0,

				"isActive": false,

				"startEvent": nil,
				"endEvent":   nil,
				"moveEvent":  nil,
			},
			Attrs: map[string]string{
				"class": config.Class + "-wrap",
			},
			Hooks: gas.Hooks{
				Mounted: func(this *gas.C) error {
					var _boundary *dom.Element
					if config.Boundary != "" {
						_boundary = dom.Doc.QuerySelector("." + config.Boundary)
						if _boundary == nil {
							dom.ConsoleError("boundary is undefined")
						}
					}

					moveEvent := event(func(event dom.Event) {
						if !this.Get("isActive").(bool) {
							return
						}

						event.PreventDefault()

						var x, y int
						if event.Type() == "touchmove" {
							t := event.JSValue().Get("touches").Get("0")
							x = t.Get("clientX").Int()
							y = t.Get("clientY").Int()
						} else {
							x = event.JSValue().Get("clientX").Int()
							y = event.JSValue().Get("clientY").Int()
						}

						if _boundary != nil {
							rect := _boundary.JSValue().Call("getBoundingClientRect")

							var (
								left   = rect.Get("left").Int()
								top    = rect.Get("top").Int()
								bottom = rect.Get("bottom").Int()
								right  = rect.Get("right").Int()

								cursorOffsetLeft   = this.Get("cursorOffsetLeft").(int)
								cursorOffsetTop    = this.Get("cursorOffsetTop").(int)
								cursorOffsetRight  = this.Get("cursorOffsetRight").(int)
								cursorOffsetBottom = this.Get("cursorOffsetBottom").(int)
							)

							if (x - cursorOffsetLeft) <= left {
								x = left + cursorOffsetLeft
							} else if (x + cursorOffsetRight) >= right {
								x = right - cursorOffsetRight
							}

							if (y - cursorOffsetTop) <= top {
								y = top + cursorOffsetTop
							} else if (y + cursorOffsetBottom) >= bottom {
								y = bottom - cursorOffsetBottom
							}
						}

						x = (x - this.Get("initialX").(int)) * config.ifXDisabled
						y = (y - this.Get("initialY").(int)) * config.ifYDisabled

						this.Set(map[string]interface{}{
							"offsetX": x,
							"offsetY": y,
						})

						if config.OnMove != nil {
							err := config.OnMove(web.ToUniteObject(event), x, y)
							if err != nil {
								this.ConsoleError(err.Error())
							}
						}
					})

					startEvent := event(func(event dom.Event) {
						if config.Handle == "" {
							_target := event.Target()
							if _target.GetAttribute("data-i").String() != childUUID && !this.Element().(*dom.Element).Contains(_target) {
								return
							}
						} else if !event.Target().ClassList().Contains(config.Handle) {
							return
						}

						var clientX, clientY int
						if event.Type() == "touchstart" {
							t := event.JSValue().Get("touches").Get("0")
							clientX = t.Get("clientX").Int()
							clientY = t.Get("clientY").Int()
						} else {
							clientX = event.JSValue().Get("clientX").Int()
							clientY = event.JSValue().Get("clientY").Int()
						}

						x := clientX - this.Get("offsetX").(int)
						y := clientY - this.Get("offsetY").(int)

						if config.OnStart != nil {
							block, err := config.OnStart(web.ToUniteObject(event), x, y)
							if err != nil {
								this.ConsoleError(err.Error())
								return
							}
							if block {
								return
							}
						}

						updatesMap := map[string]interface{}{
							"initialX": x,
							"initialY": y,
							"isActive": true,
						}

						if _boundary != nil {
							rect := dom.Doc.QuerySelector("[data-i='" + childUUID + "']").JSValue().Call("getBoundingClientRect")
							updatesMap["cursorOffsetLeft"] = clientX - rect.Get("left").Int()
							updatesMap["cursorOffsetTop"] = clientY - rect.Get("top").Int()
							updatesMap["cursorOffsetRight"] = rect.Get("right").Int() - clientX
							updatesMap["cursorOffsetBottom"] = rect.Get("bottom").Int() - clientY
						}

						this.Set(updatesMap)
					})

					endEvent := event(func(event dom.Event) {
						if !this.Get("isActive").(bool) {
							return
						}

						if config.OnEnd != nil {
							reset, err := config.OnEnd(web.ToUniteObject(event), this.Get("offsetX").(int), this.Get("offsetY").(int))
							if err != nil {
								this.ConsoleError(err.Error())
								return
							}

							if reset {
								this.Set(map[string]interface{}{
									"initialX": 0,
									"initialY": 0,

									"offsetX":  0,
									"offsetY":  0,
									"isActive": false,
								})
								return
							}
						}

						this.Set(map[string]interface{}{
							"initialX": this.Get("offsetX"),
							"initialY": this.Get("offsetY"),
							"isActive": false,
						})
					})

					addEvent(dom.Doc, "mousemove", moveEvent)
					addEvent(dom.Doc, "mousedown", startEvent)
					addEvent(dom.Doc, "mouseup", endEvent)

					this.Set(map[string]interface{}{
						"moveEvent":  moveEvent,
						"startEvent": startEvent,
						"endEvent":   endEvent,
					})

					return nil
				},
				BeforeDestroy: func(this *gas.C) error {
					if _, ok := this.Get("moveEvent").(js.Func); !ok {
						return errors.New("invalid component Data")
					}

					removeEvent(dom.Doc, "mousemove", this.Get("moveEvent").(js.Func))
					removeEvent(dom.Doc, "mousedown", this.Get("startEvent").(js.Func))
					removeEvent(dom.Doc, "mouseup", this.Get("endEvent").(js.Func))

					return nil
				},
			},
		},
		func(this *gas.Component) []interface{} {
			return []interface{}{
				gas.NE(
					&gas.C{
						UUID: childUUID,
						Binds: map[string]gas.Bind{
							"style": func() string {
								return fmt.Sprintf("transform: translate3d(%dpx, %dpx, 0px)", this.Get("offsetX"), this.Get("offsetY"))
							},
							"class": func() string {
								isActive := this.Get("isActive").(bool)
								var isActiveClass string
								if isActive {
									isActiveClass = config.Class + "-active"
								}
								return config.Class + " " + isActiveClass
							},
						},
						Hooks: gas.Hooks{
							Mounted: func(p *gas.C) error {
								_el := p.Element().(*dom.Element)

								addEvent(_el, "touchstart", this.Get("startEvent").(js.Func))
								addEvent(_el, "touchend", this.Get("endEvent").(js.Func))
								addEvent(_el, "touchcancel", this.Get("endEvent").(js.Func))
								addEvent(_el, "touchmove", this.Get("moveEvent").(js.Func))

								return nil
							},
						},
					},
					e.Body...),
			}
		})
}

func addEvent(e dom.Node, typ string, h js.Func) {
	e.JSValue().Call("addEventListener", typ, h)
}

func removeEvent(e dom.Node, typ string, h js.Func) {
	e.JSValue().Call("removeEventListener", typ, h)
}

func event(f func(event dom.Event)) js.Func {
	return js.NewEventCallback(func(v js.Value) {
		f(dom.ConvertEvent(v))
	})
}
