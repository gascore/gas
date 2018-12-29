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
			&gas.C{
				Data: map[string]interface{}{
					"show": true,
				},
				Attrs: map[string]string{
					"id": "if",
				},
			},
			func(this *gas.C) []interface{} {
				return gas.ToGetComponentList(
					gas.NE(
						&gas.C{
							Handlers: map[string]gas.Handler {
								"click": func(c *gas.C, e gas.HandlerEvent) {
									gas.WarnError(this.SetData("show", !this.GetData("show").(bool)))
								},
							},
							Tag: "button",
							Attrs: map[string]string{
								"id": "if__button",
							},
						},
						gas.NE(
							&gas.C{
								Directives:gas.Directives{
									If: func(p *gas.C) bool {
										return this.GetData("show").(bool)
									},
								},
							},
							"Show text"),
						gas.NE(
							&gas.C{
								Directives:gas.Directives{
									If: func(p *gas.C) bool {
										return !this.GetData("show").(bool)
									},
								},
							},
							"Hide text"),),
					gas.NE(
						&gas.C{
							Directives: gas.Directives{
								// If `Directives.Show == false` set `display: none` to element styles
								Show: func(c *gas.C) bool {
									return !this.GetData("show").(bool)
								},
							},
							Tag: "i",
						},
						"Hidden text",),
					gas.NE(
						&gas.C{
							Directives: gas.Directives{
								If: func(c *gas.C) bool {
									return this.GetData("show").(bool)
								},
							},
							Tag: "b",
						},
						"Public text",),)
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
