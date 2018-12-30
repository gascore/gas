package main

import (
	"fmt"
	"github.com/gascore/gas"
	"github.com/gascore/gas-web"
	"github.com/gascore/gas-web/gojs"
	"github.com/gascore/gas-web/wasm"
)

// Example application #6
//
// 'model-directive' shows how you can use component.Directive.Model
func main() {
	app, err :=
		gas.New(
			gas_web.GetBackEnd(wasm.GetDomBackEnd()),
			"app",
			&gas.Component{},
			func(main *gas.Component) []interface{} {
				return gas.ToGetComponentList(
					gas.NC(
						&gas.Component{
							Data: map[string]interface{}{
								"foo": "",
							},
							Tag: "div",
							Attrs: map[string]string{
								"id": "model__text",
								"style": "border: 1px solid #dedede; margin-bottom: 8px; padding: 4px 16px;",
							},
						},
						func (this *gas.Component) []interface{} {
							return gas.ToGetComponentList(
								fmt.Sprintf("Your text: %s", this.GetData("foo").(string)),
								gas.NE(&gas.Component{Tag: "br"}),
								gas.NE(
									&gas.Component{
									Directives: gas.Directives{
										Model: gas.ModelDirective{
											Data: "foo",
											Component: this,
										},
									},
									Tag: "input",
								},
								),
							)
						},
					),
					gas.NC(
						&gas.Component{
							Data:
							map[string]interface{}{
								"foo": "",
							},
							Tag: "div",
							Attrs:
							map[string]string{
								"id": "model__color",
								"style": "border: 1px solid #dedede; margin-bottom: 8px; padding: 4px 16px;",
							},
						},
						func (this *gas.Component) []interface{} {
							return gas.ToGetComponentList(
								"Your color: ",
								gas.NE(
									&gas.Component{
										Tag: "span",
										Binds: map[string]gas.Bind{
											"style": func(this2 *gas.C) string {
												return fmt.Sprintf("color: %s", this.GetData("foo"))
											},
										},
									},
									this.GetData("foo").(string),
								),
								gas.NE(&gas.Component{Tag: "br"}),
								gas.NE(
									&gas.Component{
										Directives: gas.Directives{
											Model: gas.ModelDirective{
												Data:      "foo",
												Component: this,
											},
										},
										Tag: "input",
										Attrs: map[string]string{
											"type": "color",
										},
									},
								),
							)
						},
					),
					gas.NC(
						&gas.Component{
							Data: map[string]interface{}{
								"foo": int(0),
							},
							Tag: "div",
							Attrs: map[string]string{
								"id": "model__range",
								"style": "border: 1px solid #dedede; margin-bottom: 8px; padding: 4px 16px;",
							},
						},
						func (this *gas.Component) []interface{} {
							return gas.ToGetComponentList(
								fmt.Sprintf("Your range: %d", this.GetData("foo").(int)),
								gas.NE(&gas.Component{Tag: "br"}),
								gas.NE(
									&gas.Component{
										Directives: gas.Directives{
											Model: gas.ModelDirective{
												Data:      "foo",
												Component: this,
											},
										},
										Tag: "input",
										Attrs: map[string]string{
											"type": "range",
										},
									},
								),
							)
						},
					),
					gas.NC(
						&gas.Component{
							Data: map[string]interface{}{
								"foo": false,
							},
							Tag: "div",
							Attrs: map[string]string{
								"id": "model__date",
								"style": "border: 1px solid #dedede; margin-bottom: 8px; padding: 4px 16px;",
							},
						},
						func (this *gas.Component) []interface{} {
							return gas.ToGetComponentList(
								fmt.Sprintf("Your checkbox: %t", this.GetData("foo").(bool)),
								gas.NE(&gas.Component{Tag: "br"}),
								gas.NE(
									&gas.Component{
										Directives: gas.Directives{
											Model: gas.ModelDirective{
												Data:      "foo",
												Component: this,
											},
										},
										Tag: "input",
										Attrs: map[string]string{
											"type": "checkbox",
										},
									},
								),
							)
						},
					),
				)
			},)
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
