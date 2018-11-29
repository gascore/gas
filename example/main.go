package main

import (
	"fmt"
	"github.com/Sinicablyat/dom"
	"github.com/Sinicablyat/gas"
)

// This all seems very weired, BUT! this will auto generate from nice-look .gas components.
// HelloWorld component will look like this:
//
//  <h1 id="hello-world" class="greeting h1">
// 		Hello world!
//  </h1>
//  <i id="italiano" class="greeting">
//		Ciao mondo!
//  </i>
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
				gas.NilData,
				gas.NilMethods,
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
				gas.NilData,
				gas.NilMethods,
				gas.NilBinds,
				gas.NilHandlers,
				"i",
				map[string]string{
					"id":    "italiano",
					"class": "greeting",
					"style": "margin-right: 12px;",
				},
				func(p gas.Component) interface{} {
					return "Ciao mondo!" // I'm not italian, but i love films about mafia
				})
			},
			func(p gas.Component) interface{} {
				return gas.NewComponent(
					&p,
					map[string]interface{}{
						"click": 0,
					},
					gas.NilData,
					gas.NilMethods,
					gas.NilBinds,
					map[string]gas.Handler {
						"click": func(c gas.Component, e dom.Event) {
							currentClick := c.GetData("click").(int)
							_ = c.SetData("click", currentClick+1)
							dom.ConsoleDir(fmt.Sprintf("before: %d, after: %d", currentClick, c.GetData("click")))
						},
					},
					"button",
					map[string]string{
						"id": "clicker",
					},
					func(p gas.Component) interface{} {
						return fmt.Sprintf("Yout clicked me %d", p.GetData("click").(int))
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
