package main

import (
	"fmt"
	"github.com/Sinicablyat/gas"
	"github.com/Sinicablyat/wasm"
)

// Example application #10
//
// 'binds' shows how you can use component.Binds
func main() {
	app, err :=
		gas.New(
			wasm.GetBackEnd(),
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

											// REMEMBER!
											// Bind attributes appends over normal attributes, you will lose your normal attribute value
											return fmt.Sprintf("%s; background-color: rgb(%d, %d, %d)", this2.Attrs["style"], foo, 255-foo, foo)
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
										return "// color: rgb(x, 255-x, x)"
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
									"min": "0",
									"max": "255",
								},
							})
					},)
			},
			)
	must(err)

	err = gas.Init(app)
	must(err)
	gas.KeepAlive()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
