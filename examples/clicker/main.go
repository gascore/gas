package main

import (
	"fmt"
	"github.com/Sinicablyat/dom"
	"github.com/Sinicablyat/gas"
)

// Example application #2
//
// 'clicker' shows how you can add handlers and change component.Data
func main() {
	app, err :=
		gas.New(
			"app",
			func(p gas.Component) interface{} {
				return gas.NewComponent(
					&p,
					map[string]interface{}{
						"click": 0,
					},
					gas.NilMethods,
					gas.NilDirectives,
					gas.NilBinds,
					gas.NilHandlers,
					"h1",
					map[string]string{
						"id": "clicker",
					},
					func(this gas.Component) interface{} {
						return gas.NewComponent(
							&p,
							gas.NilData,
							gas.NilMethods,
							gas.NilDirectives,
							gas.NilBinds,
							map[string]gas.Handler {
								"click": func(c gas.Component, e dom.Event) {
									currentClick := this.GetData("click").(int)
									gas.WarnError(this.SetData("click", currentClick+1))
								},
							},
							"button",
							map[string]string{
								"id": "clicker__button", // I love BEM
							},
							func(this2 gas.Component) interface{} {
								return "Click me!"
							})
					},
					func(this gas.Component) interface{} {
						return fmt.Sprintf("You clicked button %d times", this.GetData("click").(int))
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
