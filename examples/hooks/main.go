package main

import (
	"fmt"
	"github.com/gascore/gas"
	"github.com/gascore/gas-web"
)

// Example application #9
//
// 'hooks' shows how you can use component.Hooks
func main() {
	app, err :=
		gas.New(
			web.GetBackEnd(),
			"app",
			&gas.C{
				Data: map[string]interface{}{
					"show":    true,
					"counter": 0,
				},
				Attrs: map[string]string{
					"id": "hooks",
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
								"id": "hooks__button",
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
							If: func(c *gas.C) bool {
								return this.Get("show").(bool)
							},
							Hooks: gas.Hooks{
								Mounted: func(this2 *gas.C) error {
									this.ConsoleLog("Visible text is created!")
									this.SetValue("counter", this.Get("counter").(int)+1)
									return nil
								},
								BeforeDestroy: func(this2 *gas.C) error {
									this.ConsoleLog("Visible text will destroy!")
									return nil
								},
							},
							Tag: "i",
						},
						fmt.Sprintf("Visible text"),
					),
					gas.NE(
						&gas.C{
							Else: true,
							Hooks: gas.Hooks{
								Mounted: func(this2 *gas.C) error {
									this.ConsoleLog("Hidden text is created!")
									this.SetValue("counter", this.Get("counter").(int)+1)
									return nil
								},
								BeforeDestroy: func(this2 *gas.C) error {
									this.ConsoleLog("Hidden text will destroy!")
									return nil
								},
							},
							Tag: "i",
						},
						fmt.Sprintf("Hidden text (you show hidden text %d times)", this.Get("counter")),
					),
				)
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
