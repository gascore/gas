package main

import (
	"github.com/gascore/gas"
	"github.com/gascore/gas/web"
)

// Example application #3
//
// 'if-directive' shows how you can use component.Directive.If
func main() {
	app, err :=
		gas.New(
			web.GetBackEnd(),
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
				return gas.CL(
					gas.NE(
						&gas.C{
							Handlers: map[string]gas.Handler{
								"click": func(c *gas.C, e gas.Object) {
									this.SetValue("show", !this.Get("show").(bool))
								},
							},
							Tag: "button",
							Attrs: map[string]string{
								"id": "if__button",
							},
						},
						gas.NE(
							&gas.C{
								If: func(p *gas.C) bool {
									return this.Get("show").(bool)
								},
							},
							"Show text"),
						gas.NE(
							&gas.C{
								If: func(p *gas.C) bool {
									return !this.Get("show").(bool)
								},
							},
							"Hide text")),
					gas.NE(
						&gas.C{
							// If `Directives.Show == false` set `display: none` to element styles
							Show: func(c *gas.C) bool {
								return !this.Get("show").(bool)
							},
							Tag: "i",
						},
						"Hidden text"),
					gas.NE(
						&gas.C{
							If: func(c *gas.C) bool {
								return this.Get("show").(bool)
							},
							Tag: "b",
						},
						"Public text"))
			})
	must(err)

	err = gas.Init(app)
	must(err)
	web.KeepAlive()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
