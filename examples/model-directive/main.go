package main

import (
	"fmt"
	"github.com/gascore/gas"
	"github.com/gascore/gas-web"
)

// Example application #6
//
// 'model-directive' shows how you can use component.Directive.Model
func main() {
	app, err :=
		gas.New(
			web.GetBackEnd(),
			"app",
			&gas.Component{},
			func(main *gas.Component) []interface{} {
				return gas.CL(
					gas.NC(
						&gas.Component{
							Data: map[string]interface{}{
								"foo": "",
							},
							Tag: "div",
							Attrs: map[string]string{
								"id":    "model__text",
								"style": "border: 1px solid #dedede; margin-bottom: 8px; padding: 4px 16px;",
							},
						},
						func(this *gas.Component) []interface{} {
							return gas.CL(
								fmt.Sprintf("Your text: %s", this.Get("foo").(string)),
								gas.NE(&gas.Component{Tag: "br"}),
								gas.NE(
									&gas.Component{
										Model: gas.ModelDirective{
											Data:      "foo",
											Component: this,
										},
										Tag: "input",
									},
								),
							)
						},
					),
					gas.NC(
						&gas.Component{
							Data: map[string]interface{}{
								"foo": "",
							},
							Tag: "div",
							Attrs: map[string]string{
								"id":    "model__color",
								"style": "border: 1px solid #dedede; margin-bottom: 8px; padding: 4px 16px;",
							},
						},
						func(this *gas.Component) []interface{} {
							return gas.CL(
								"Your color: ",
								gas.NE(
									&gas.Component{
										Tag: "span",
										Binds: map[string]gas.Bind{
											"style": func() string {
												return fmt.Sprintf("color: %s", this.Get("foo"))
											},
										},
									},
									this.Get("foo").(string),
								),
								gas.NE(&gas.Component{Tag: "br"}),
								gas.NE(
									&gas.Component{
										Model: gas.ModelDirective{
											Data:      "foo",
											Component: this,
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
								"id":    "model__range",
								"style": "border: 1px solid #dedede; margin-bottom: 8px; padding: 4px 16px;",
							},
						},
						func(this *gas.Component) []interface{} {
							return gas.CL(
								fmt.Sprintf("Your range: %d", this.Get("foo").(int)),
								gas.NE(&gas.Component{Tag: "br"}),
								gas.NE(
									&gas.Component{
										Model: gas.ModelDirective{
											Data:      "foo",
											Component: this,
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
								"id":    "model__date",
								"style": "border: 1px solid #dedede; margin-bottom: 8px; padding: 4px 16px;",
							},
						},
						func(this *gas.Component) []interface{} {
							return gas.CL(
								fmt.Sprintf("Your checkbox: %t", this.Get("foo").(bool)),
								gas.NE(&gas.Component{Tag: "br"}),
								gas.NE(
									&gas.Component{
										Model: gas.ModelDirective{
											Data:      "foo",
											Component: this,
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
			})
	must(err)

	err = gas.Init(app)
	must(err)
	web.KeepAlive()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
