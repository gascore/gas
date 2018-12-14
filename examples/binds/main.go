package main

import (
	"fmt"
	"github.com/Sinicablyat/gas"
)

// Example application #10
//
// 'binds' shows how you can use component.Binds
func main() {
	app, err :=
		gas.New(
			"app",
			func(p *gas.Component) interface{} {
				return gas.NewComponent(
					&gas.Component{
						ParentC: p,
						Data: map[string]interface{}{
							"foo": int(0),
						},
						Tag: "div",
						Attrs: map[string]string{
							"id": "model__range",
							"style": "border: 1px solid #dedede; margin-bottom: 8px; padding: 4px 16px;",
						},
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(
							&gas.Component{
								Tag: "div",
								Attrs: map[string]string{
									"style": "display: flex;",
								},
							},
							func(this2 *gas.Component) interface{} {
								foo, ok := this.GetData("foo").(int)
								gas.WarnIfNot(ok)
								return fmt.Sprintf("Your range: %d", foo)
							},
							func(this2 *gas.Component) interface{} {
								return gas.NewComponent(&gas.Component{
									ParentC: this,
									Binds: map[string]gas.Bind{
										"style": func(this2 *gas.Component) string {
											foo, ok := this.GetData("foo").(int)
											gas.WarnIfNot(ok)

											var color string
											switch {
											case foo%10 == 0:
												color = "background-color: red;"
											case foo%4 == 0:
												color = "background-color: green;"
											case foo%2 == 0:
												color = "background-color: blue;"
											default:
												color = "background-color: purple;"
											}

											// REMEMBER!
											// Bind attributes appends over normal attributes, you will lose your normal attribute value
											return fmt.Sprintf("%s; %s", this2.Attrs["style"], color)
										},
									},
									Attrs: map[string]string{
										"style": "width: 48px; height: 36px; margin: 0 18px; border-radius: 4px;",
									},
									Tag: "div",
								},)
							},
							func(this2 *gas.Component) interface{} {
								return gas.NewComponent(
									&gas.Component{
										Attrs: map[string]string{
											"style": "color: darkgray;",
										},
										Tag: "i",
									},
									func(this3 *gas.Component) interface{} {
										return "// if the number is a multiple of 10 - red, 4 - green, 2 - blue, other - purple"
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
								Directives: gas.Directives{
									Model: gas.ModelDirective{
										Data: "foo",
										Component: this,
									},
								},
								Tag: "input",
								Attrs: map[string]string{
									"type": "range",
								},
							})
					},)
			},
			)
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
