package main

import (
	"fmt"
	"github.com/Sinicablyat/gas"
	"github.com/Sinicablyat/gas/core"
)

// Example application #6
//
// 'model-directive' shows how you can use component.Directive.Model
func main() {
	app, err :=
		gas.NewWasm(
			"app",
			func(p *core.Component) interface{} {
				return core.NewComponent(
					&core.Component{
						ParentC: p,
						Data: map[string]interface{}{
							"foo": "",
						},
						Tag: "div",
						Attrs: map[string]string{
							"id": "model__text",
							"style": "border: 1px solid #dedede; margin-bottom: 8px; padding: 4px 16px;",
						},
					},
					func(this *core.Component) interface{} {
						foo, ok := this.GetData("foo").(string)
						gas.WarnIfNot(ok)
						return fmt.Sprintf("Your text: %s", foo)
					},
					func(this *core.Component) interface{} {
						return core.NewComponent(&core.Component{ParentC: this, Tag: "br"})
					},
					func(this *core.Component) interface{} {
						return core.NewComponent(
							&core.Component{
								ParentC: this,
								Directives: core.Directives{
									Model: core.ModelDirective{
										Data: "foo",
										Component: this,
									},
								},
								Tag: "input",
							})
					},)
			},
			func(p *core.Component) interface{} {
				return core.NewComponent(
					&core.Component{
						ParentC: p,
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
					func(this *core.Component) interface{} {
						foo, ok := this.GetData("foo").(string)
						gas.WarnIfNot(ok)
						return fmt.Sprintf("Your color: %s", foo)
					},
					func(this *core.Component) interface{} {
						return core.NewComponent(&core.Component{ParentC: this, Tag: "br"})
					},
					func(this *core.Component) interface{} {
						return core.NewComponent(
							&core.Component{
								ParentC: this,
								Directives: core.Directives{
									Model: core.ModelDirective{
										Data: "foo",
										Component: this,
									},
								},
								Tag: "input",
								Attrs: map[string]string{
									"type": "color",
								},
							})
					},)
			},
			func(p *core.Component) interface{} {
				return core.NewComponent(
					&core.Component{
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
					func(this *core.Component) interface{} {
						foo, ok := this.GetData("foo").(int)
						gas.WarnIfNot(ok)
						return fmt.Sprintf("Your range: %d", foo)
					},
					func(this *core.Component) interface{} {
						return core.NewComponent(&core.Component{ParentC: this, Tag: "br"})
					},
					func(this *core.Component) interface{} {
						return core.NewComponent(
							&core.Component{
								ParentC: this,
								Directives: core.Directives{
									Model: core.ModelDirective{
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
			func(p *core.Component) interface{} {
				return core.NewComponent(
					&core.Component{
						ParentC: p,
						Data: map[string]interface{}{
							"foo": false,
						},
						Tag: "div",
						Attrs: map[string]string{
							"id": "model__date",
							"style": "border: 1px solid #dedede; margin-bottom: 8px; padding: 4px 16px;",
						},
					},
					func(this *core.Component) interface{} {
						foo, ok := this.GetData("foo").(bool)
						gas.WarnIfNot(ok)
						return fmt.Sprintf("Your checkbox: %t", foo)
					},
					func(this *core.Component) interface{} {
						return core.NewComponent(&core.Component{ParentC: this, Tag: "br"})
					},
					func(this *core.Component) interface{} {
						return core.NewComponent(
							&core.Component{
								ParentC: this,
								Directives: core.Directives{
									Model: core.ModelDirective{
										Data: "foo",
										Component: this,
									},
								},
								Tag: "input",
								Attrs: map[string]string{
									"type": "checkbox",
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
