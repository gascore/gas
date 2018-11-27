package main

import (
	"github.com/Sinicablyat/gas"
)

// This all seems very weired, BUT! this will auto generate from nice-look .gas components.
// HelloWorld component will look like this:
//
//  <h1 id=hello-world>
// 		Hello world!
//  </h1>
//  <i id="italiano">
//		Ciao mondo!
//  </i>
func main() {
	app, err :=
		gas.New(
			"app",
			func(p gas.Component) interface{} {
				return gas.NewComponent(map[string]interface{}{
					"hello": "Hello world!",
				}, gas.NilData, "h1", "hello-world", []string{"greeting", "h1"}, gas.NilAttrs).
					AddBinds(gas.NilBinds).
					AddChildes(
						func(this gas.Component) interface{} {
							return this.GetData("hello")
						})
			},
			func(p gas.Component) interface{} {
				return gas.NewComponent(gas.NilData, gas.NilData, "i", "italinao", []string{"greeting"}, gas.NilAttrs).
					AddBinds(gas.NilBinds).
					AddChildes(
						func(p gas.Component) interface{} {
							return "Ciao mondo!" // I'm not italian, but i love pizza
						})
			})
	must(err)

	err = app.Init()
	must(err)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
