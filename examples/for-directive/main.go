package main

import (
	"fmt"
	"github.com/Sinicablyat/gas"
	"github.com/Sinicablyat/wasm"
)

// Example application #4
//
// 'if-directive' shows how you can use component.Directive.For
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
							"arr": []interface{}{"click", "here", "if you want to see some magic"},
						},
						Tag: "ul",
						Attrs: map[string]string{
							"id": "list",
						},
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(
							&gas.Component{
								ParentC: this,
								Directives: gas.Directives{
									For: gas.ForDirective{
										Data: "arr",
										Render: func(i int, el interface{}, this *gas.Component) []gas.GetComponent {
											return gas.ToGetComponentList(
												func(this2 *gas.Component) interface{} {
													return fmt.Sprintf("%d: %s", i+1, el)
												},)
										},
										Component: this,
									},
								},
								Handlers: map[string]gas.Handler {
									"click": func(c *gas.Component, e interface{}) {
										arr := this.GetData("arr").([]interface{})
										arr = append(arr, "Hello!") // hello, Annoy-o-Tron
										gas.WarnError(this.SetData("arr", arr))
									},
								},
								Tag: "li",
							}) // In components with FOR Directive childes are ignored
					},
					func(this *gas.Component) interface{} {
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
