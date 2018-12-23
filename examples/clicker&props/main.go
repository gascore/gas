package main

import (
	"github.com/Sinicablyat/gas"
	"github.com/Sinicablyat/gas-web"
	"github.com/Sinicablyat/gas-web/wasm"
)

// Example application #2
//
// 'clicker&props' shows how you can add handlers, change component.Data and use external components
func main() {
	app, err :=
		gas.New(
			gas_web.GetBackEnd(wasm.GetDomBackEnd()),
			"app",
			func(p *gas.Component) interface{} {
				return gas.NewComponent(
					&gas.Component{
						ParentC:p,
						Data: map[string]interface{}{
							"click": 0,
						},
						Methods: map[string]gas.Method{
							"addClick": func(this *gas.Component, i ...interface{}) error {
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
					func(this *gas.Component) interface{} {
						return gas.NewComponent(
							&gas.Component{
								ParentC: this,
								Handlers: map[string]gas.Handler {
									"click.left": func(this2 *gas.Component, e interface{}) {
										gas.WarnError(this.Method("addClick"))
									},
									// you need to click button once (for target it)
									"keyup.control": func(this2 *gas.Component, e interface{}) {
										gas.WarnError(this.Method("addClick"))
									},
									"keyup.a": func(this2 *gas.Component, e interface{}) {
										gas.WarnError(this.Method("addClick"))
									},
									"keyup.s": func(this2 *gas.Component, e interface{}) {
										gas.WarnError(this.Method("addClick"))
									},
									"keyup.d": func(this2 *gas.Component, e interface{}) {
										gas.WarnError(this.Method("addClick"))
									},
									"keyup.f": func(this2 *gas.Component, e interface{}) {
										gas.WarnError(this.Method("addClick"))
									},
								},
								Tag: "button",
								Attrs: map[string]string{
									"id": "clicker__button", // I love BEM
								},
							},
							func(this2 *gas.Component) interface{} {
								return "Click me!"
							})
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(
							&gas.Component{
								ParentC: this,
								Tag: "span",
								Attrs: map[string]string{
									"id": "needful_wrapper",
								},
							},
							func(this2 *gas.Component) interface{} {
								return "You clicked button: "
							},
							func(this2 *gas.Component) interface{} {
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
