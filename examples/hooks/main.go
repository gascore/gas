package main

import (
	"fmt"
	"github.com/Sinicablyat/dom"
	"github.com/gascore/gas"
	"github.com/gascore/gas-web"
	"github.com/gascore/gas-web/wasm"
)

// Example application #9
//
// 'hooks' shows how you can use component.Hooks
func main() {
	app, err :=
		gas.New(
			gas_web.GetBackEnd(wasm.GetDomBackEnd()),
			"app",
			&gas.C{
				Data: map[string]interface{}{
					"show":    true,
					"counter": 0,
				},
				Hooks: gas.Hooks{
					BeforeCreate: func(this *gas.C) error {
						dom.ConsoleLog("Component is being created!")
						return nil
					},
				},
				Attrs: map[string]string{
					"id": "hooks",
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
								"id": "hooks__button",
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
								Directives:
								gas.Directives{
									If: func(c *gas.C) bool {
										return !this.GetData("show").(bool)
									},
								},
								Hooks:
								gas.Hooks{
									Created: func(this2 *gas.C) error {
										dom.ConsoleLog("Hidden text is created!")
										return this.SetData("counter", this.GetData("counter").(int)+1)
									},
									Destroyed: func(this2 *gas.C) error {
										dom.ConsoleLog("Hidden text was destroyed!")
										return nil
									},
								},
								Tag: "i",
							},
							fmt.Sprintf("Hidden text (you show hidden text %d times)", this.GetData("counter")),
						),
					)
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
