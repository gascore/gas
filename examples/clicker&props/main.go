package main

import (
	"github.com/Sinicablyat/dom"
	"github.com/Sinicablyat/gas"
	"github.com/Sinicablyat/gas/core"
)

// Example application #2
//
// 'clicker&props' shows how you can add handlers, change component.Data and use external components
func main() {
	app, err :=
		gas.NewWasm(
			"app",
			func(p *core.Component) interface{} {
				return core.NewComponent(
					&core.Component{
						ParentC:p,
						Data: map[string]interface{}{
							"click": 0,
						},
						Methods: map[string]core.Method{
							"addClick": func(this *core.Component, i ...interface{}) error {
								currentClick := this.GetData("click").(int)
								gas.WarnError(this.SetData("click", currentClick+1))
								return nil
							},
						},
						Tag: "h1",
						Attrs: map[string]string{
							"id": "clicker&props",
						},
					},
					func(this *core.Component) interface{} {
						return core.NewComponent(
							&core.Component{
								ParentC: this,
								Handlers: map[string]core.Handler {
									"click.left": func(this2 *core.Component, e dom.Event) {
										gas.WarnError(this.Method("addClick"))
									},
									// you need to click button once (for target it)
									"keyup.control": func(this2 *core.Component, e dom.Event) {
										gas.WarnError(this.Method("addClick"))
									},
									"keyup.a": func(this2 *core.Component, e dom.Event) {
										gas.WarnError(this.Method("addClick"))
									},
									"keyup.s": func(this2 *core.Component, e dom.Event) {
										gas.WarnError(this.Method("addClick"))
									},
									"keyup.d": func(this2 *core.Component, e dom.Event) {
										gas.WarnError(this.Method("addClick"))
									},
									"keyup.f": func(this2 *core.Component, e dom.Event) {
										gas.WarnError(this.Method("addClick"))
									},
								},
								Tag: "button",
								Attrs: map[string]string{
									"id": "clicker__button", // I love BEM
								},
							},
							func(this2 *core.Component) interface{} {
								return "Click me!"
							})
					},
					func(this *core.Component) interface{} {
						return core.NewComponent(
							&core.Component{
								ParentC: this,
								Tag: "span",
								Attrs: map[string]string{
									"id": "needful_wrapper",
								},
							},
							func(this2 *core.Component) interface{} {
								return "You clicked button: "
							},
							func(this2 *core.Component) interface{} {
								// It's EXTERNAL component!
								return GetNumberViewer(this, this.GetData("click").(int))
							})
					})
			},)
	must(err)

	err = gas.Init(app)
	must(err)
	gas.KeepAlive()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
