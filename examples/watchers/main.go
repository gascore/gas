package main

import (
	"fmt"
	"github.com/gascore/gas"
	"github.com/gascore/gas-web"
	"github.com/gascore/gas-web/wasm"
)

// Example application #8
//
// 'watchers' shows how you can use component.Watchers
func main() {
	app, err :=
		gas.New(
			gas_web.GetBackEnd(wasm.GetDomBackEnd()),
			"app",
			&gas.Component{
				Data: map[string]interface{}{
					"show":               true,
					"watcherIsTriggered": false,
				},
				Watchers: map[string]gas.Watcher{
					"show": func(this *gas.Component, new interface{}, old interface{}) error {
						this.ConsoleLog(fmt.Sprintf("Watcher is triggered! New value: %t, old value: %t", new, old))

						err := this.SetDataFree("watcherIsTriggered", true)
						if err != nil {
							this.WarnError(err)
							return err
						}

						return nil
					},
				},
				Attrs: map[string]string{
					"id": "if",
				},
			},
			func(this *gas.Component) []interface{} {
				return gas.CL(
					gas.NE(
						&gas.Component{
							Handlers: map[string]gas.Handler{
								"click": func(c *gas.Component, e gas.Object) {
									this.WarnError(this.SetData("show", !this.GetData("show").(bool)))
								},
							},
							Tag: "button",
							Attrs: map[string]string{
								"id": "if__button",
							},
						},
						gas.NE(
							&gas.C{
								Directives: gas.Directives{
									Show: func(c *gas.C) bool {
										return this.GetData("show").(bool)
									},
								},
								Tag: "i",
							},
							"Show text"),
						gas.NE(
							&gas.C{
								Directives: gas.Directives{
									Show: func(c *gas.C) bool {
										return !this.GetData("show").(bool)
									},
								},
								Tag: "i",
							},
							"Hidden text",
						),
					),
					gas.NE(
						&gas.Component{
							Directives: gas.Directives{
								If: func(c *gas.Component) bool {
									return !this.GetData("show").(bool)
								},
							},
							Tag: "i",
						},
						"Hidden text",
					),
					gas.NE(
						&gas.Component{
							Directives: gas.Directives{
								If: func(this2 *gas.Component) bool {
									watcherIsTriggered, ok := this.GetData("watcherIsTriggered").(bool)
									this.WarnIfNot(ok)
									return watcherIsTriggered
								},
							},
							Tag: "strong",
							Attrs: map[string]string{
								"style": "color: red;",
							},
						},
						"Watcher is triggered!",
					),
				)
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
