package main

import (
	"github.com/Sinicablyat/dom"
	"github.com/Sinicablyat/gas"
)

// Example application #2
//
// 'clicker&props' shows how you can add handlers, change component.Data and use external components
func main() {
	app, err :=
		gas.New(
			"app",
			func(p *gas.Component) interface{} {
				return gas.NewComponent(
					p,
					map[string]interface{}{
						"click": 0,
					},
					gas.NilMethods,
					gas.NilDirectives,
					gas.NilBinds,
					gas.NilHandlers,
					"h1",
					map[string]string{
						"id": "clicker&props",
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(
							p,
							gas.NilData,
							gas.NilMethods,
							gas.NilDirectives,
							gas.NilBinds,
							map[string]gas.Handler {
								"click": func(c gas.Component, e dom.Event) {
									currentClick := this.GetData("click").(int)
									gas.WarnError(this.SetData("click", currentClick+1))
								},
							},
							"button",
							map[string]string{
								"id": "clicker__button", // I love BEM
							},
							func(this2 *gas.Component) interface{} {
								return "Click me!"
							})
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(
							this,
							gas.NilData,
							gas.NilMethods,
							gas.NilDirectives,
							gas.NilBinds,
							gas.NilHandlers,
							"span",
							map[string]string{
								"id": "needful_wrapper",
							},
							func(this2 *gas.Component) interface{} {
								return "You clicked button: "
							},
							func(this2 *gas.Component) interface{} {
								// It's EXTERNAL component!
								return GetNumberViewer(this2, this.GetData("click").(int))
							})
					})
			},)
	must(err)

	err = app.Init()
	must(err)
	gas.KeepAlive()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
