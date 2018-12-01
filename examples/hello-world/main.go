package main

import (
	"github.com/Sinicablyat/gas"
)

// Example application #1
//
// 'hello-world' shows how you can create components, component.Data and component.Attributes
func main() {
	app, err :=
		gas.New(
			"app",
			func(p gas.Component) interface{} {
			return gas.NewComponent(
				&p,
				map[string]interface{}{
					"hello": "Hello world!",
				},
				gas.NilMethods,
				gas.NilDirectives,
				gas.NilBinds,
				gas.NilHandlers,
				"h1",
				map[string]string{
					"id":    "hello-world",
					"class": "greeting h1",
				},
				func(this gas.Component) interface{} {
					return this.GetData("hello")
				})
			},
			func(p gas.Component) interface{} {
			return gas.NewComponent(
				&p,
				gas.NilData,
				gas.NilMethods,
				gas.NilDirectives,
				gas.NilBinds,
				gas.NilHandlers,
				"i",
				map[string]string{
					"id":    "italiano",
					"class": "greeting",
					"style": "margin-right: 12px;",
				},
				func(this gas.Component) interface{} {
					return "Ciao mondo!" // I'm not italian, but i love films about mafia
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
