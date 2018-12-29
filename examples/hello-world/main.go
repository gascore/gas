package main

import (
	"github.com/Sinicablyat/gas"
	"github.com/Sinicablyat/gas-web"
	"github.com/Sinicablyat/gas-web/wasm"
)

// Example application #1
//
// 'hello-world' shows how you can create components, component.Data and component.Attributes
func main() {
	app, err :=
		gas.New(
			//gas_web.GetBackEnd(gojs.GetDomBackEnd()),
			gas_web.GetBackEnd(wasm.GetDomBackEnd()),
			"app",
			&gas.Component{
				Data: map[string]interface{}{
					"hello": "Hello world!",
				},
			},
			func(this *gas.Component) []interface{} {
				return gas.ToGetComponentList(
					gas.NewBasicComponent(
						&gas.Component{
							Tag: "h1",
							Attrs: map[string]string{
								"id":    "hello-world",
								"class": "greeting h1",
							},
						},
						this.GetData("hello")),
					gas.NewBasicComponent(
						&gas.Component{
							Tag: "i",
							Attrs: map[string]string{
								"id":    "italiano",
								"class": "greeting",
								"style": "margin-right: 12px;",
							},
						},
						"Ciao mondo!"))
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
