package main

import (
	"fmt"
	"github.com/Sinicablyat/dom"
	"github.com/Sinicablyat/gas"
)

// Example application #5
//
// 'methods&computed' shows how you can use component.Methods and component.Computed.
func main() {
	app, err := // as you can see i love 'clicker' example
		gas.New(
			"app",
			func(p *gas.Component) interface{} {
				return gas.NewComponent(
					&gas.Component{
						ParentC: p,
						Data: map[string]interface{}{
							"show": true,
							"number": 1,
						},
						// What the difference between Methods and Computed?
						// Methods will do business things.
						// Computed will return value from data, libraries, /dev/random, e.t.c. with some changes (or just raw)
						Methods: map[string]gas.Method{
							//	Component updates after method is ended.
							//	Therefore we can change data millions times, but it will updates once.
							//	And u should know it
							"toggle": func(this *gas.Component, values ...interface{}) error {
								// in methods we call only SetDataFree because we don't need updates component after changes
								_ = this.SetDataFree("show", !this.GetData("show").(bool))
								_ = this.SetDataFree("show", !this.GetData("show").(bool))
								_ = this.SetDataFree("show", !this.GetData("show").(bool))
								_ = this.SetDataFree("show", !this.GetData("show").(bool))
								_ = this.SetDataFree("show", !this.GetData("show").(bool))

								if this.GetData("show").(bool) {
									_ = this.SetDataFree("number", this.GetData("number").(int)+1)
								}

								return nil
							},
						},
						Computeds: map[string]gas.Computed{
							"number": func(this *gas.Component, values ...interface{}) (interface{}, error) {
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
					func(this *gas.Component) interface{} {
						// For pass method or computed to sub component you need to get a *pocket* version.
						// Pocket method/computed can be executed in sub component, in sub sub component, in e.t.c.
						pocketToggle, err := this.GetPocketMethod("toggle")
						gas.WarnError(err)

						return getButton(this, pocketToggle)
					},
					func(this *gas.Component) interface{} {
						pocketNumber, err := this.GetPocketComputed("number")
						gas.WarnError(err)

						return getHiddenText(this, this.GetData("show").(bool), pocketNumber)
					})
			},)
	must(err)

	err = app.Init()
	must(err)
	gas.KeepAlive()
}

func getButton(this *gas.Component, toggleMethod gas.PocketMethod) *gas.Component {
	return gas.NewComponent(
		&gas.Component{
			ParentC: this,
			Handlers: map[string]gas.Handler {
				"click": func(c *gas.Component, e dom.Event) {
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
		func(this2 *gas.Component) interface{} {
			if this.GetData("show").(bool) {
				return "Show text"
			} else {
				return "Hide text"
			}
		})
}

func getHiddenText(this *gas.Component, isShow bool, getNumber gas.PocketComputed) *gas.Component {
	return gas.NewComponent(
		&gas.Component{
			ParentC: this,
			Directives: gas.Directives{
				If: func(c *gas.Component) bool {
					return !isShow
				},
			},
			Tag: "i",
		},
		func(this2 *gas.Component) interface{} {
			return "Hidden text"
		},
		func (this2 *gas.Component) interface{} {
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
