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
					p,
					map[string]interface{}{
						"show": true,
						"number": 1,
					},
					// What the difference between Methods and Computed?
					// Methods will do business things.
					// Computed will return value from data, libraries, /dev/random, e.t.c. with some changes (or just raw)
					map[string]gas.Method{
						"toggle": func(this *gas.Component) error {
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
					map[string]gas.Computed{
						"number": func(this *gas.Component) (interface{}, error) {
							currentNumber, ok := this.GetData("number").(int)
							gas.WarnIfNot(ok) // it's good practise to your data for valid type
							explanation := fmt.Sprintf("You showed hidden text: %d times", currentNumber)
							return explanation, nil
						},
					},
					gas.NilDirectives,
					gas.NilBinds,
					gas.NilHandlers,
					"h1",
					map[string]string{
						"id": "M&C",
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(
							p,
							gas.NilData,
							gas.NilMethods,
							gas.NilComputeds,
							gas.NilDirectives,
							gas.NilBinds,
							map[string]gas.Handler {
								"click": func(c gas.Component, e dom.Event) {
									/*
										Component updates after method is ended.
										Therefore we can change data millions times, but it will updates once.
										And u should know it
									*/
									gas.WarnError(this.Method("toggle"))
								},
							},
							"button",
							map[string]string{
								"id": "M&C__button",
							},
							func(this2 *gas.Component) interface{} {
								if this.GetData("show").(bool) {
									return "Show text"
								} else {
									return "Hide text"
								}
							})
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(
							p,
							gas.NilData,
							gas.NilMethods,
							gas.NilComputeds,
							gas.Directives{
								If: func(c *gas.Component) bool {
									return !this.GetData("show").(bool)
								},
							},
							gas.NilBinds,
							gas.NilHandlers,
							"i",
							gas.NilAttrs,
							func(this2 *gas.Component) interface{} {
								return "Hidden text"
							},
							func (this2 *gas.Component) interface{} {
								value, err := this.Computed("number")
								gas.WarnError(err) // always check for error
								return fmt.Sprintf("  (%s)", value)
							})
					},)
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
