package main

import (
	"github.com/Sinicablyat/dom"
	"github.com/Sinicablyat/gas"
)

// Example application #3
//
// 'if-directive' shows how you can use component.Directive.If
func main() {
	app, err :=
		gas.New(
			"app",
			func(p *gas.Component) interface{} {
				return gas.NewComponent(
					&gas.Component{
						ParentC: p,
						Data: map[string]interface{}{
							"show": true,
						},
						Tag: "h1",
						Attrs: map[string]string{
							"id": "if",
						},
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(
							&gas.Component{
								ParentC: this,
								Handlers: map[string]gas.Handler {
									"click": func(c *gas.Component, e dom.Event) {
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
								ParentC: this,
								Directives: gas.Directives{
									If: func(c *gas.Component) bool {
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
								ParentC: this,
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
