package main

import (
	"github.com/Sinicablyat/gas"
	"github.com/Sinicablyat/gas/core"
)

// Example application #1
//
// 'hello-world' shows how you can create components, component.Data and component.Attributes
func main() {
	app, err :=
		gas.NewWasm(
			"app",
			func(p *core.Component) interface{} {
			return core.NewComponent(
				&core.Component{
					ParentC: p,
					Data:
						map[string]interface{}{
							"hello": "Hello world!",
						},
					Tag: "h1",
					Attrs:
						map[string]string{
							"id":    "hello-world",
							"class": "greeting h1",
						},
				},
				func(this *core.Component) interface{} {
					return this.GetData("hello")
				})
			},
			func(p *core.Component) interface{} {
			return core.NewComponent(
				&core.Component{
					ParentC:p,
					Tag: "i",
					Attrs:
						map[string]string{
							"id":    "italiano",
							"class": "greeting",
							"style": "margin-right: 12px;",
						},
				},
				func(this *core.Component) interface{} {
					return "Ciao mondo!" // I'm not italian, but i love films about mafia
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
