package main

import (
	"fmt"
	"github.com/Sinicablyat/dom"
	"github.com/Sinicablyat/gas"
)

// Example application #4
//
// 'if-directive' shows how you can use component.Directive.For
func main() {
	app, err :=
		gas.New(
			"app",
			func(p *gas.Component) interface{} {
				return gas.NewComponent(
					p,
					map[string]interface{}{
						"arr": []interface{}{"click", "here", "if you want to see some magic"},
					},
					gas.NilWatchers,
					gas.NilMethods,
					gas.NilComputeds,
					gas.NilDirectives,
					gas.NilBinds,
					gas.NilHandlers,
					"ul",
					map[string]string{
						"id": "list",
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(
							this,
							gas.NilData,
							gas.NilWatchers,
							gas.NilMethods,
							gas.NilComputeds,
							gas.Directives{
								If: gas.NilIfDirective,
								For: gas.ForDirective{
									Data: "arr",
									Render: func(i int, el interface{}, this *gas.Component) []gas.GetComponent {
										return gas.ToGetComponentList(
											func(this2 *gas.Component) interface{} {
												return fmt.Sprintf("%d: %s", i+1, el)
											},)
									},
								},
								HTML: gas.NilHTMLDirective,
							},
							gas.NilBinds,
							map[string]gas.Handler {
								"click": func(c *gas.Component, e dom.Event) {
									arr := this.GetData("arr").([]interface{})
									arr = append(arr, "Hello!") // hello, Annoy-o-Tron
									gas.WarnError(this.SetData("arr", arr))
								},
							},
							"li",
							gas.NilAttrs,) // In components with FOR Directive childes are ignored
					},
					func(this *gas.Component) interface{} {
						return "end of list"
					})
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
