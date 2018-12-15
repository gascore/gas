package main

import (
	"fmt"
	"github.com/Sinicablyat/dom"
	"github.com/Sinicablyat/gas"
)

// Example application #11
//
// 'todo' shows how you how to build basic TODO example
func main() {
	app, err :=
		gas.New(
			"app",
			func(p *gas.Component) interface{} {
				return gas.NewComponent(
					&gas.Component{
						ParentC: p,
						Data: map[string]interface{}{
							"currentText": "",
							"list": []interface{}{},
						},
						Tag: "div",
						Attrs: map[string]string{
							"id": "todo",
						},
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(
							&gas.Component{
								ParentC: this,
								Tag: "div",
							},
							func(this2 *gas.Component) interface{} {
								return gas.NewComponent(
									&gas.Component{
										ParentC: this2,
										Directives: gas.Directives{
											Model: gas.ModelDirective{
												Data: "currentText",
												Component: this,
											},
										},
										Tag: "input",
									})
							},
							func(this2 *gas.Component) interface{} {
								return gas.NewComponent(
									&gas.Component{
										ParentC: this2,
										Handlers: map[string]gas.Handler{
											"click": func(this3 *gas.Component, e dom.Event) {
												list := this.GetData("list").([]interface{})
												currentText := this.GetData("currentText").(string)

												if len(currentText) == 0 {
													return
												}

												list  = append(list, currentText)

												gas.WarnError(this.SetData("list", list))
												gas.WarnError(this.SetData("currentText", ""))
											},
										},
										Attrs: map[string]string{
											"type": "button",
										},
										Tag: "button",
									},
									func(this3 *gas.Component) interface{} {
										return "Add"
									})
							},)
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(&gas.Component{ParentC: this, Tag: "br"})
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(
							&gas.Component{
								ParentC: this,
								Tag: "ul",
							},
							func(this2 *gas.Component) interface{} {
								return gas.NewComponent(
									&gas.Component{
										ParentC: this2,
										Tag: "li",
										Directives:gas.Directives{
											For:gas.ForDirective{
												Data: "list",
												Component: this,
												Render: func(i int, value interface{}, this3 *gas.Component) []gas.GetComponent {
													return gas.ToGetComponentList(
														func(this4 *gas.Component) interface{} {
															return gas.NewComponent(
																&gas.Component{
																	Tag: "i",
																	Attrs: map[string]string{
																		"style": "margin-right: 8px",
																	},
																},
																func(this5 *gas.Component) interface{} {
																	return fmt.Sprintf("%d: %s", i+1, value)
																})
														},
														func(this4 *gas.Component) interface{} {
															return gas.NewComponent(
																&gas.Component{
																	Tag: "button",
																	Handlers: map[string]gas.Handler{
																		"click": func(this5 *gas.Component, e dom.Event) {
																			list := this.GetData("list").([]interface{})

																			gas.WarnError(this.DoWithUpdate(func() error {
																				return this.SetData("list", remove(list, i))
																			}))

																			gas.WarnError(this.SetData("currentText", ""))
																		},
																	},
																},
																func(this5 *gas.Component) interface{} {
																	return "Remove"
																})
														})
												},
											},
										},
									})
							})
					},)
			},
		)
	must(err)

	err = app.Init()
	must(err)
	gas.KeepAlive()
}

func remove(a []interface{}, i int) []interface{} {
	copy(a[i:], a[i+1:]) // Shift a[i+1:] left one index
	a[len(a)-1] = ""     // Erase last element (write zero value)
	a = a[:len(a)-1]     // Truncate slice

	return a
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
