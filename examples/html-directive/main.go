package main

import (
	"github.com/gascore/gas"
	"github.com/gascore/gas-web"
)

// Example application #7
//
// 'html-directive' shows how you can use component.Directive.HTML
func main() {
	app, err :=
		gas.New(
			web.GetBackEnd(),
			"app",
			&gas.C{
				Data: map[string]interface{}{
					"articleText": `
<h1>
	Lorem ipsum dolor sit amet, consectetur adipiscing elit.
</h1>
<p>
	Vivamus arcu nibh, sodales nec lectus ut, vestibulum porta est. Nunc in odio eu tellus feugiat volutpat vitae a erat.
</p>
<p>
	<i>Phasellus sit amet suscipit urna</i>. 
	Quisque vitae risus lobortis, aliquam orci at, pulvinar urna. Quisque vitae lobortis libero.
	Nullam a faucibus dolor. Ut eu turpis et purus mollis ullamcorper. Vivamus interdum felis quis volutpat volutpat. Mauris id auctor nisi.
</p>
<hr/>
<p>
	<strong>Integer aliquam tellus nunc, ac dapibus felis pulvinar viverra</strong>. 
Donec dapibus dolor in massa vehicula ornare. Duis molestie velit vitae purus consectetur pulvinar. Aliquam ac purus placerat, laoreet tortor at, aliquet ex.
</p>
<h3>
	Nulla facilisi. Donec mattis auctor finibus.
</h3>`,
					"helloText":       `<h1>To see article click button!</h1>`,
					"isArticleActive": false,
				},
			},
			func(this *gas.C) []interface{} { // don't use childes if you have v-html
				return gas.CL(
					gas.NE(
						&gas.C{
							Handlers: map[string]gas.Handler{
								"click": func(this2 *gas.C, e gas.Object) {
									currentIsArticleActive := this.Get("isArticleActive").(bool)
									this.WarnError(this.SetValue("isArticleActive", !currentIsArticleActive))
								},
							},
							Tag: "button",
						},
						gas.NE(
							&gas.C{
								Directives: gas.Directives{
									If: func(p *gas.C) bool {
										return this.Get("isArticleActive").(bool)
									},
								},
							},
							"Hide article"),
						gas.NE(
							&gas.C{
								Directives: gas.Directives{
									If: func(p *gas.C) bool {
										return !this.Get("isArticleActive").(bool)
									},
								},
							},
							"Show article"),
					),
					gas.NE(
						&gas.C{
							Directives: gas.Directives{
								HTML: gas.HTMLDirective{Render: func(this2 *gas.C) string {
									isArticleActive, ok := this.Get("isArticleActive").(bool)

									var html string
									if isArticleActive {
										html, ok = this.Get("articleText").(string)
									} else {
										html, ok = this.Get("helloText").(string)
									}
									this.WarnIfNot(ok)

									return html
								}},
							},
							Tag: "article",
							Attrs: map[string]string{
								"id":    "article",
								"style": `border: 1px solid #dedede;padding: 2px 4px;margin-top:12px;`,
							},
						}))
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
