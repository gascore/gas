package main

import (
	"fmt"
	"github.com/Sinicablyat/dom"
	"github.com/Sinicablyat/gas"
	"github.com/Sinicablyat/gas/core"
)

// Example application #9
//
// 'hooks' shows how you can use component.Hooks
func main() {
	app, err :=
		gas.NewWasm(
			"app",
			func(p *core.Component) interface{} {
				return core.NewComponent(
					&core.Component{
						ParentC: p,
						Data: map[string]interface{}{
							"show":    true,
							"counter": 0,
						},
						Hooks: core.Hooks{
							BeforeCreate: func(this *core.Component) error {
								dom.ConsoleLog("Component is being created!")
								return nil
							},
						},
						Tag: "h1",
						Attrs: map[string]string{
							"id": "hooks",
						},
					},
					func(this *core.Component) interface{} {
						return core.NewComponent(
							&core.Component{
								ParentC: this,
								Handlers: map[string]core.Handler {
									"click": func(c *core.Component, e dom.Event) {
										gas.WarnError(this.SetData("show", !this.GetData("show").(bool)))
									},
								},
								Tag: "button",
								Attrs: map[string]string{
									"id": "hooks__button",
								},
							},
							func(this2 *core.Component) interface{} {
								if this.GetData("show").(bool) {
									return "Show text"
								} else {
									return "Hide text"
								}
							})
					},
					func(this *core.Component) interface{} {
						return core.NewComponent(
							&core.Component{
								ParentC: this,
								Directives:
									core.Directives{
										If: func(c *core.Component) bool {
											return !this.GetData("show").(bool)
										},
									},
								Hooks:
									core.Hooks{
										Created: func(this2 *core.Component) error {
											dom.ConsoleLog("Hidden text is created!")
											return this.SetData("counter", this.GetData("counter").(int)+1)
										},
										Destroyed: func(this2 *core.Component) error {
											dom.ConsoleLog("Hidden text was destroyed!")
											return nil
										},
									},
								Tag: "i",
							},
							func(this2 *core.Component) interface{} {
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
