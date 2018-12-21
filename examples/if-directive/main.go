package main

import (
	"github.com/Sinicablyat/dom"
	"github.com/Sinicablyat/gas"
	"github.com/Sinicablyat/gas/core"
)

// Example application #3
//
// 'if-directive' shows how you can use component.Directive.If
func main() {
	app, err :=
		gas.NewWasm(
			"app",
			func(p *core.Component) interface{} {
				return core.NewComponent(
					&core.Component{
						ParentC: p,
						Data: map[string]interface{}{
							"show": true,
						},
						Tag: "h1",
						Attrs: map[string]string{
							"id": "if",
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
									"id": "if__button",
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
								Directives: core.Directives{
									// If `Directives.Show == false` set `display: none` to element styles
									Show: func(c *core.Component) bool {
										return !this.GetData("show").(bool)
									},
								},
								Tag: "i",
							},
							func(this2 *core.Component) interface{} {
								return "Hidden text"
							})
					},
					func(this *core.Component) interface{} {
						return core.NewComponent(
							&core.Component{
								ParentC: this,
								Directives: core.Directives{
									If: func(c *core.Component) bool {
										return this.GetData("show").(bool)
									},
								},
								Tag: "b",
							},
							func(this2 *core.Component) interface{} {
								return "Public text"
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
