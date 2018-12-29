package main

import (
	"fmt"
	"github.com/Sinicablyat/dom"
	"github.com/Sinicablyat/gas"
	"github.com/Sinicablyat/gas-web"
	"github.com/Sinicablyat/gas-web/wasm"
)

// Example application #5
//
// 'methods&computed' shows how you can use component.Methods and component.Computed.
func main() {
	app, err :=
		gas.New(
			gas_web.GetBackEnd(wasm.GetDomBackEnd()),
			"app",
			&gas.Component{
				Data: map[string]interface{}{
					"show": true,
					"number": 1,
				},
				// What the difference between Methods and Computed?
				// Methods will do business things.
				// Computed will return value from data, libraries, /dev/random, e.t.c. with some changes (or just raw)
				Methods: map[string]gas.Method{
					"toggle": func(this *gas.Component, values ...interface{}) error {
						_ = this.SetData("show", !this.GetData("show").(bool))

						if this.GetData("show").(bool) {
							_ = this.SetData("number", this.GetData("number").(int)+1)
						}

						return nil
					},
				},
				// Computeds can be cached
				Computeds: map[string]gas.Computed{
					"number": func(this *gas.Component, values ...interface{}) (interface{}, error) {
						dom.ConsoleLog(fmt.Sprintf("Some values: %s", values[0].(string)))

						currentNumber, ok := this.GetData("number").(int)
						gas.WarnIfNot(ok) // it's good practise to your data for valid type
						explanation := fmt.Sprintf("You showed hidden text: %d times", currentNumber)
						return explanation, nil
					},
				},
				Attrs: map[string]string{
					"id": "M&C",
				},
			},
			func(this *gas.Component) []interface{} {
				return gas.ToGetComponentList(
					getButton(this.GetData("show").(bool), this.GetPocketMethod("toggle")),
					getHiddenText(this.GetData("show").(bool), this.GetPocketComputed("number")),)
			},)
	must(err)

	err = gas.Init(app)
	must(err)
	gas.KeepAlive()
}

func getButton(show bool, toggleMethod gas.PocketMethod) *gas.Component {
	return gas.NE(
		&gas.Component{
			Handlers: map[string]gas.Handler {
				"click": func(c *gas.Component, e gas.HandlerEvent) {
					// Of course we can use method for `this`.
					// But if we want to pass method to child not from `this` we need to pass a pocket method/computed.
					gas.WarnError(toggleMethod())
				},
			},
			Tag: "button",
			Attrs:
				map[string]string{
					"id": "M&C__button",
				},
		},
		gas.NE(
			&gas.C{
				Directives:gas.Directives{
					If: func(p *gas.C) bool {
						return show
					},
				},
			},
			"Show text"),
		gas.NE(
			&gas.C{
				Directives:gas.Directives{
					If: func(p *gas.C) bool {
						return !show
					},
				},
			},
			"Hide text"))
}

func getHiddenText(show bool, getNumber gas.PocketComputed) *gas.Component {
	return gas.NE(
		&gas.Component{
			Directives: gas.Directives{
				If: func(c *gas.Component) bool {
					return !show
				},
			},
			Tag: "i",
		},
		"Hidden text",
		fmt.Sprintf("  (%s)", getNumber("something for computed")))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
