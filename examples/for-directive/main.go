package main

import (
	"fmt"
	"github.com/Sinicablyat/gas"
)

// Example application #3
//
// 'if-directive' shows how you can use component.Directives
func main() {
	app, err :=
		gas.New(
			"app",
			func(p gas.Component) interface{} {
				return gas.NewComponent(
					&p,
					map[string]interface{}{
						"arr": []interface{}{"one", "two", "three"},
					},
					gas.NilMethods,
					gas.NilDirectives,
					gas.NilBinds,
					gas.NilHandlers,
					"h1",
					map[string]string{
						"id": "if",
					},
					func(this gas.Component) interface{} {
						return gas.NewComponent(
							&p,
							gas.NilData,
							gas.NilMethods,
							gas.Directives{
								If: gas.NilIfDirective,
								For: gas.ForDirective{
									Data: "arr",
									Render: func(arr []interface{}) []interface{} {
										var elements []interface{}
										for i, el := range arr {
											elements = append(elements, fmt.Sprintf("%d: %s", i, el))
										}
										return elements
									},
								},
							},
							gas.NilBinds,
							gas.NilHandlers,
							"i",
							gas.NilAttrs) // In components with For Directive childes are ignored
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
