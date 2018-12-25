package main

import (
	"fmt"
	"github.com/Sinicablyat/dom"
	"github.com/Sinicablyat/gas"
	"github.com/Sinicablyat/gas-web"
	"github.com/Sinicablyat/gas-web/wasm"
)

// Example application #9
//
// 'hooks' shows how you can use component.Hooks
func main() {
	app, err :=
		gas.New(
			gas_web.GetBackEnd(wasm.GetDomBackEnd()),
			"app",
			func(p *gas.Component) interface{} {
				return gas.NewComponent(
					&gas.Component{
						ParentC: p,
						Data: map[string]interface{}{
							"show":    true,
							"counter": 0,
						},
						Hooks: gas.Hooks{
							BeforeCreate: func(this *gas.Component) error {
								dom.ConsoleLog("Component is being created!")
								return nil
							},
						},
						Tag: "h1",
						Attrs: map[string]string{
							"id": "hooks",
						},
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(
							&gas.Component{
								ParentC: this,
								Handlers: map[string]gas.Handler {
									"click": func(c *gas.Component, e gas.HandlerEvent) {
										gas.WarnError(this.SetData("show", !this.GetData("show").(bool)))
									},
								},
								Tag: "button",
								Attrs: map[string]string{
									"id": "hooks__button",
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
								ParentC: this,
								Directives:
									gas.Directives{
										If: func(c *gas.Component) bool {
											return !this.GetData("show").(bool)
										},
									},
								Hooks:
									gas.Hooks{
										Created: func(this2 *gas.Component) error {
											dom.ConsoleLog("Hidden text is created!")
											return this.SetData("counter", this.GetData("counter").(int)+1)
										},
										Destroyed: func(this2 *gas.Component) error {
											dom.ConsoleLog("Hidden text was destroyed!")
											return nil
										},
									},
								Tag: "i",
							},
							func(this2 *gas.Component) interface{} {
								return "Hidden text " + fmt.Sprintf("(you show hidden text %d times)", this.GetData("counter"))
							},)
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
