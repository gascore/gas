package main

import (
	"fmt"
	"github.com/Sinicablyat/dom"
	"github.com/Sinicablyat/gas"
	"github.com/Sinicablyat/gas/core"
)

// Example application #5
//
// 'methods&computed' shows how you can use component.Methods and component.Computed.
func main() {
	app, err := // as you can see i love 'clicker' example
		gas.NewWasm(
			"app",
			func(p *core.Component) interface{} {
				return core.NewComponent(
					&core.Component{
						ParentC: p,
						Data: map[string]interface{}{
							"show": true,
							"number": 1,
						},
						// What the difference between Methods and Computed?
						// Methods will do business things.
						// Computed will return value from data, libraries, /dev/random, e.t.c. with some changes (or just raw)
						Methods: map[string]core.Method{
							"toggle": func(this *core.Component, values ...interface{}) error {
								_ = this.SetData("show", !this.GetData("show").(bool))

								if this.GetData("show").(bool) {
									_ = this.SetData("number", this.GetData("number").(int)+1)
								}

								return nil
							},
						},
						// Computeds can be cached
						Computeds: map[string]core.Computed{
							"number": func(this *core.Component, values ...interface{}) (interface{}, error) {
								dom.ConsoleLog(fmt.Sprintf("Some values: %s", values[0].(string)))

								currentNumber, ok := this.GetData("number").(int)
								gas.WarnIfNot(ok) // it's good practise to your data for valid type
								explanation := fmt.Sprintf("You showed hidden text: %d times", currentNumber)
								return explanation, nil
							},
						},
						Tag: "h1",
						Attrs: map[string]string{
							"id": "M&C",
						},
					},
					func(this *core.Component) interface{} {
						// For pass method or computed to sub component you need to get a *pocket* version.
						// Pocket method/computed can be executed in sub component, in sub sub component, in e.t.c.
						pocketToggle, err := this.GetPocketMethod("toggle")
						gas.WarnError(err)

						return getButton(this, pocketToggle)
					},
					func(this *core.Component) interface{} {
						pocketNumber, err := this.GetPocketComputed("number")
						gas.WarnError(err)

						return getHiddenText(this, this.GetData("show").(bool), pocketNumber)
					})
			},)
	must(err)

	err = gas.Init(app)
	must(err)
	gas.KeepAlive()
}

func getButton(this *core.Component, toggleMethod core.PocketMethod) *core.Component {
	return core.NewComponent(
		&core.Component{
			ParentC: this,
			Handlers: map[string]core.Handler {
				"click": func(c *core.Component, e dom.Event) {
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
		func(this2 *core.Component) interface{} {
			if this.GetData("show").(bool) {
				return "Show text"
			} else {
				return "Hide text"
			}
		})
}

func getHiddenText(this *core.Component, isShow bool, getNumber core.PocketComputed) *core.Component {
	return core.NewComponent(
		&core.Component{
			ParentC: this,
			Directives: core.Directives{
				If: func(c *core.Component) bool {
					return !isShow
				},
			},
			Tag: "i",
		},
		func(this2 *core.Component) interface{} {
			return "Hidden text"
		},
		func (this2 *core.Component) interface{} {
			value, err := getNumber("something for computed")
			gas.WarnError(err) // always check for error
			return fmt.Sprintf("  (%s)", value)
		})
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
