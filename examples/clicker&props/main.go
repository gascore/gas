package main

import (
	"github.com/gascore/gas"
	"github.com/gascore/gas-web"
)

// Example application #2
//
// 'clicker&props' shows how you can add handlers, change component.Data and use external components
func main() {
	app, err :=
		gas.New(
			web.GetBackEnd(),
			"app",
			&gas.C{
				Data: map[string]interface{}{
					"click": 0,
				},
				Methods: map[string]gas.Method{
					"addClick": func(this *gas.C, i ...interface{}) (interface{}, error) {
						currentClick := this.Get("click").(int)

						err := this.SetValue("click", currentClick+1)
						if err != nil {
							return nil, err
						}

						return nil, nil
					},
				},
				Attrs: map[string]string{
					"id": "clicker&props",
				},
			},
			func(this *gas.C) []interface{} {
				return gas.CL(
					gas.NE(
						&gas.C{
							Handlers: map[string]gas.Handler{
								"click.left": func(this2 *gas.C, e gas.Object) {
									this.Method("addClick")
								},
								// you need to click button once (for target it)
								"keyup.control": func(this2 *gas.C, e gas.Object) {
									this.Method("addClick")
								},
								"keyup.a": func(this2 *gas.C, e gas.Object) {
									this.Method("addClick")
								},
								"keyup.s": func(this2 *gas.C, e gas.Object) {
									this.Method("addClick")
								},
								"keyup.d": func(this2 *gas.C, e gas.Object) {
									this.Method("addClick")
								},
								"keyup.f": func(this2 *gas.C, e gas.Object) {
									this.Method("addClick")
								},
							},
							Tag: "button",
							Attrs: map[string]string{
								"id": "clicker__button", // I love BEM
							},
						},
						"Click me!"),
					gas.NE(
						&gas.C{
							Tag: "i",
							Attrs: map[string]string{
								"id": "needful_wrapper",
							},
						},
						"You clicked button: ",
						GetNumberViewer(this.Get("click").(int))))
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
