package main

import (
	"fmt"
	"github.com/Sinicablyat/dom"
	"github.com/Sinicablyat/gas"
	"github.com/Sinicablyat/gas/core"
)

// Example application #4
//
// 'if-directive' shows how you can use component.Directive.For
func main() {
	app, err :=
		gas.NewWasm(
			"app",
			func(p *core.Component) interface{} {
				return core.NewComponent(
					&core.Component{
						ParentC: p,
						Data: map[string]interface{}{
							"arr": []interface{}{"click", "here", "if you want to see some magic"},
						},
						Tag: "ul",
						Attrs: map[string]string{
							"id": "list",
						},
					},
					func(this *core.Component) interface{} {
						return core.NewComponent(
							&core.Component{
								ParentC: this,
								Directives: core.Directives{
									For: core.ForDirective{
										Data: "arr",
										Render: func(i int, el interface{}, this *core.Component) []core.GetComponent {
											return gas.ToGetComponentList(
												func(this2 *core.Component) interface{} {
													return fmt.Sprintf("%d: %s", i+1, el)
												},)
										},
										Component: this,
									},
								},
								Handlers: map[string]core.Handler {
									"click": func(c *core.Component, e dom.Event) {
										arr := this.GetData("arr").([]interface{})
										arr = append(arr, "Hello!") // hello, Annoy-o-Tron
										gas.WarnError(this.SetData("arr", arr))
									},
								},
								Tag: "li",
							}) // In components with FOR Directive childes are ignored
					},
					func(this *core.Component) interface{} {
						return "end of list"
					})
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
