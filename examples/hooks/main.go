package main

import (
	"github.com/Sinicablyat/dom"
	"github.com/Sinicablyat/gas"
)

// Example application #9
//
// 'hooks' shows how you can use component.Hooks
func main() {
	app, err :=
		gas.New(
			"app",
			func(p *gas.Component) interface{} {
				return gas.NewComponent(
					p,
					map[string]interface{}{
						"show": true,
					},
					gas.NilWatchers,
					gas.NilMethods,
					gas.NilComputeds,
					gas.NilDirectives,
					gas.NilBinds,
					gas.Hooks{
						BeforeCreate: func(this *gas.Component) error {
							dom.ConsoleLog("Component is being created!")
							return nil
						},
					},
					gas.NilHandlers,
					"h1",
					map[string]string{
						"id": "hooks",
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(
							this,
							gas.NilData,
							gas.NilWatchers,
							gas.NilMethods,
							gas.NilComputeds,
							gas.NilDirectives,
							gas.NilBinds,
							gas.NilHooks,
							map[string]gas.Handler {
								"click": func(c *gas.Component, e dom.Event) {
									gas.WarnError(this.SetData("show", !this.GetData("show").(bool)))
								},
							},
							"button",
							map[string]string{
								"id": "hooks__button",
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
							this,
							gas.NilData,
							gas.NilWatchers,
							gas.NilMethods,
							gas.NilComputeds,
							gas.Directives{
								If: func(c *gas.Component) bool {
									return !this.GetData("show").(bool)
								},
								HTML: gas.NilHTMLDirective,
							},
							gas.NilBinds,
							gas.Hooks{
								Created: func(this2 *gas.Component) error {
									dom.ConsoleLog("Hidden text is created!")
									return nil
								},
								Destroyed: func(this2 *gas.Component) error {
									dom.ConsoleLog("Hidden text was destroyed!")
									return nil
								},
							},
							gas.NilHandlers,
							"i",
							gas.NilAttrs,
							func(this2 *gas.Component) interface{} {
								return "Hidden text"
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
