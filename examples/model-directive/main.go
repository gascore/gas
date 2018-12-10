package main

import (
	"fmt"
	"github.com/Sinicablyat/gas"
)

// Example application #6
//
// 'model-directive' shows how you can use component.Directive.Model
func main() {
	app, err :=
		gas.New(
			"app",
			func(p *gas.Component) interface{} {
				return gas.NewComponent(
					p,
					map[string]interface{}{
						"foo": "",
					},
					gas.NilMethods,
					gas.NilComputeds,
					gas.NilDirectives,
					gas.NilBinds,
					gas.NilHandlers,
					"div",
					map[string]string{
						"id": "model__text",
						"style": "border: 1px solid #dedede; margin-bottom: 8px; padding: 4px 16px;",
					},
					func(this *gas.Component) interface{} {
						foo, ok := this.GetData("foo").(string)
						gas.WarnIfNot(ok)
						return fmt.Sprintf("Your text: %s", foo)
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(this, gas.NilData, gas.NilMethods, gas.NilComputeds, gas.NilDirectives, gas.NilBinds, gas.NilHandlers, "br", gas.NilAttrs)
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(
							this,
							gas.NilData,
							gas.NilMethods,
							gas.NilComputeds,
							gas.Directives{
								If: gas.NilIfDirective,
								Model: gas.ModelDirective{
									Data: "foo",
									Component: this,
								},
								HTML: gas.NilHTMLDirective,
							},
							gas.NilBinds,
							gas.NilHandlers,
							"input",
							gas.NilAttrs)
					},)
			},
			func(p *gas.Component) interface{} {
				return gas.NewComponent(
					p,
					map[string]interface{}{
						"foo": "",
					},
					gas.NilMethods,
					gas.NilComputeds,
					gas.NilDirectives,
					gas.NilBinds,
					gas.NilHandlers,
					"div",
					map[string]string{
						"id": "model__color",
						"style": "border: 1px solid #dedede; margin-bottom: 8px; padding: 4px 16px;",
					},
					func(this *gas.Component) interface{} {
						foo, ok := this.GetData("foo").(string)
						gas.WarnIfNot(ok)
						return fmt.Sprintf("Your color: %s", foo)
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(this, gas.NilData, gas.NilMethods, gas.NilComputeds, gas.NilDirectives, gas.NilBinds, gas.NilHandlers, "br", gas.NilAttrs)
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(
							this,
							gas.NilData,
							gas.NilMethods,
							gas.NilComputeds,
							gas.Directives{
								If: gas.NilIfDirective,
								Model: gas.ModelDirective{
									Data: "foo",
									Component: this,
								},
								HTML: gas.NilHTMLDirective,
							},
							gas.NilBinds,
							gas.NilHandlers,
							"input",
							map[string]string{
								"type": "color",
							})
					},)
			},
			func(p *gas.Component) interface{} {
				return gas.NewComponent(
					p,
					map[string]interface{}{
						"foo": int(0),
					},
					gas.NilMethods,
					gas.NilComputeds,
					gas.NilDirectives,
					gas.NilBinds,
					gas.NilHandlers,
					"div",
					map[string]string{
						"id": "model__range",
						"style": "border: 1px solid #dedede; margin-bottom: 8px; padding: 4px 16px;",
					},
					func(this *gas.Component) interface{} {
						foo, ok := this.GetData("foo").(int)
						gas.WarnIfNot(ok)
						return fmt.Sprintf("Your range: %d", foo)
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(this, gas.NilData, gas.NilMethods, gas.NilComputeds, gas.NilDirectives, gas.NilBinds, gas.NilHandlers, "br", gas.NilAttrs)
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(
							this,
							gas.NilData,
							gas.NilMethods,
							gas.NilComputeds,
							gas.Directives{
								If: gas.NilIfDirective,
								Model: gas.ModelDirective{
									Data: "foo",
									Component: this,
								},
								HTML: gas.NilHTMLDirective,
							},
							gas.NilBinds,
							gas.NilHandlers,
							"input",
							map[string]string{
								"type": "range",
							})
					},)
			},
			func(p *gas.Component) interface{} {
				return gas.NewComponent(
					p,
					map[string]interface{}{
						"foo": false,
					},
					gas.NilMethods,
					gas.NilComputeds,
					gas.NilDirectives,
					gas.NilBinds,
					gas.NilHandlers,
					"div",
					map[string]string{
						"id": "model__date",
						"style": "border: 1px solid #dedede; margin-bottom: 8px; padding: 4px 16px;",
					},
					func(this *gas.Component) interface{} {
						foo, ok := this.GetData("foo").(bool)
						gas.WarnIfNot(ok)
						return fmt.Sprintf("Your checkbox: %t", foo)
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(this, gas.NilData, gas.NilMethods, gas.NilComputeds, gas.NilDirectives, gas.NilBinds, gas.NilHandlers, "br", gas.NilAttrs)
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(
							this,
							gas.NilData,
							gas.NilMethods,
							gas.NilComputeds,
							gas.Directives{
								If: gas.NilIfDirective,
								Model: gas.ModelDirective{
									Data: "foo",
									Component: this,
								},
								HTML: gas.NilHTMLDirective,
							},
							gas.NilBinds,
							gas.NilHandlers,
							"input",
							map[string]string{
								"type": "checkbox",
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
