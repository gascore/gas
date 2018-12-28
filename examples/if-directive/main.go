package main

import (
	"github.com/Sinicablyat/gas"
	"github.com/Sinicablyat/gas-web"
	"github.com/Sinicablyat/gas-web/wasm"
)

// Example application #3
//
// 'if-directive' shows how you can use component.Directive.If
func main() {
	app, err :=
		gas.New(
			gas_web.GetBackEnd(wasm.GetDomBackEnd()),
			"app",
			&gas.Component{
				Data: map[string]interface{}{
					"show": true,
				},
				Attrs: map[string]string{
					"id": "if",
				},
			},
			func(this *gas.Component) interface{} {
				return gas.NewComponent(
					&gas.Component{
						Handlers: map[string]gas.Handler {
							"click": func(c *gas.Component, e gas.HandlerEvent) {
								gas.WarnError(this.SetData("show", !this.GetData("show").(bool)))
							},
						},
						Tag: "button",
						Attrs: map[string]string{
							"id": "if__button",
						},
					},
					func(this2 *gas.Component) interface{} {
						if this.GetData("show").(bool) {
							return "Show text"
						} else {
							return "Hide text"
						}
					})
			},
			func(this *gas.Component) interface{} {
				return gas.NewComponent(
					&gas.Component{
						Directives: gas.Directives{
							// If `Directives.Show == false` set `display: none` to element styles
							Show: func(c *gas.Component) bool {
								return !this.GetData("show").(bool)
							},
						},
						Tag: "i",
					},
					func(this2 *gas.Component) interface{} {
						return "Hidden text"
					})
			},
			func(this *gas.Component) interface{} {
						return gas.NewComponent(
							&gas.Component{
								Directives: gas.Directives{
									If: func(c *gas.Component) bool {
										return this.GetData("show").(bool)
									},
								},
								Tag: "b",
							},
							func(this2 *gas.Component) interface{} {
								return "Public text"
							})
					})
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
