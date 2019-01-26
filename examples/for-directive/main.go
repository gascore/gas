package main

import (
	"fmt"
	"github.com/gascore/gas"
	"github.com/gascore/gas-web"
)

// Example application #4
//
// 'if-directive' shows how you can use component.Directive.For
func main() {
	app, err :=
		gas.New(
			gas_web.GetBackEnd(),
			"app",
			&gas.C{
				Data: map[string]interface{}{
					"arr": []interface{}{"click", "here", "if you want to see some magic"},
				},
				Tag: "ul",
				Attrs: map[string]string{
					"id": "list",
				},
			},
			func(this *gas.C) []interface{} {
				return gas.CL(
					gas.NE(
						&gas.C{
							Tag: "ul",
						},
						gas.NewFor("arr", this, func(i int, el interface{}) interface{} {
							return gas.NE(
								&gas.C{
									Handlers: map[string]gas.Handler{
										"click": func(c *gas.C, e gas.Object) {
											arr := this.GetData("arr").([]interface{})
											arr = append(arr, "Hello!") // hello, Annoy-o-Tron
											this.WarnError(this.SetData("arr", arr))
										},
									},
									Tag: "li",
								},
								fmt.Sprintf("%d: %s", i+1, el))
						}),
						"end of list"))
			})
	must(err)

	err = gas.Init(app)
	must(err)
	dom.KeepAlive()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
