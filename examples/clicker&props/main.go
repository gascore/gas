package main

import (
	"github.com/Sinicablyat/dom"
	"github.com/Sinicablyat/gas"
)

// Example application #2
//
// 'clicker&props' shows how you can add handlers, change component.Data and use external components
func main() {
	app, err :=
		gas.New(
			"app",
			func(p *gas.Component) interface{} {
				return gas.NewComponent(
					&gas.Component{
						ParentC:p,
						Data: map[string]interface{}{
							"click": 0,
						},
						Tag: "h1",
						Attrs: map[string]string{
							"id": "clicker&props",
						},
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(
							&gas.Component{
								ParentC: this,
								Handlers: map[string]gas.Handler {
									"click": func(c *gas.Component, e dom.Event) {
										currentClick := this.GetData("click").(int)
										gas.WarnError(this.SetData("click", currentClick+1))
									},
								},
								Tag: "button",
								Attrs: map[string]string{
									"id": "clicker__button", // I love BEM
								},
							},
							func(this2 *gas.Component) interface{} {
								return "Click me!"
							})
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(
							&gas.Component{
								ParentC: this,
								Tag: "span",
								Attrs: map[string]string{
									"id": "needful_wrapper",
								},
							},
							func(this2 *gas.Component) interface{} {
								return "You clicked button: "
							},
							func(this2 *gas.Component) interface{} {
								// It's EXTERNAL component!
								return GetNumberViewer(this, this.GetData("click").(int))
							})
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
