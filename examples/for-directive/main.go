package main

import (
	"fmt"
	"github.com/Sinicablyat/gas"
	"github.com/Sinicablyat/gas-web"
	"github.com/Sinicablyat/gas-web/wasm"
)

// Example application #4
//
// 'if-directive' shows how you can use component.Directive.For
func main() {
	app, err :=
		gas.New(
			gas_web.GetBackEnd(wasm.GetDomBackEnd()),
			"app",
			&gas.Component{
				Data: map[string]interface{}{
					"arr": []interface{}{"click", "here", "if you want to see some magic"},
				},
				Tag: "ul",
				Attrs: map[string]string{
					"id": "list",
				},
			},
			func (this *gas.Component) []interface{} {
				return gas.ToGetComponentList(
					gas.NewBasicComponent(
					&gas.Component{
						Tag: "ul",
					},
					gas.NewFor("arr", this, func(i int, el interface{}) interface {} {
						return gas.NewBasicComponent(
							&gas.Component{
								Handlers: map[string]gas.Handler {
									"click": func(c *gas.Component, e gas.HandlerEvent) {
										arr := this.GetData("arr").([]interface{})
										arr = append(arr, "Hello!") // hello, Annoy-o-Tron
										gas.WarnError(this.SetData("arr", arr))
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
	gas.KeepAlive()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
