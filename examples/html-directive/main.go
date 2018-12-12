package main

import (
	"github.com/Sinicablyat/dom"
	"github.com/Sinicablyat/gas"
)

// Example application #7
//
// 'html-directive' shows how you can use component.Directive.HTML
func main() {
	app, err :=
		gas.New(
			"app",
			func(p *gas.Component) interface{} {
				return gas.NewComponent(
					p,
					map[string]interface{}{
						"articleText":
						`
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
						"helloText": `<h1>To see article click button!</h1>`,
						"isArticleActive": false,
					},
					gas.NilWatchers,
					gas.NilMethods,
					gas.NilComputeds,
					gas.NilDirectives,
					gas.NilBinds,
					gas.NilHooks,
					gas.NilHandlers,
					"main",
					gas.NilAttrs,
					func(this *gas.Component) interface{} { // don't use childes if you have v-html
						return gas.NewComponent(
							this,
							gas.NilData,
							gas.NilWatchers,
							gas.NilMethods,
							gas.NilComputeds,
							gas.NilDirectives,
							gas.NilBinds,
							gas.NilHooks,
							map[string]gas.Handler{
								"click": func(this2 *gas.Component, e dom.Event) {
									currentIsArticleActive := this.GetData("isArticleActive").(bool)
									gas.WarnError(this.SetData("isArticleActive", !currentIsArticleActive))
								},
							},
							"button",
							gas.NilAttrs,
							func(this2 *gas.Component) interface{} {
								isArticleActive, ok := this.GetData("isArticleActive").(bool)
								gas.WarnIfNot(ok)

								if isArticleActive {
									return "Hide article"
								} else {
									return "Show article"
								}
							})
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
								HTML: gas.HTMLDirective{Render: func(this2 *gas.Component) string {
									isArticleActive, ok := this.GetData("isArticleActive").(bool)

									var html string
									if isArticleActive {
										html, ok = this.GetData("articleText").(string)
									} else {
										html, ok = this.GetData("helloText").(string)
									}
									gas.WarnIfNot(ok)

									return html
								},},
							},
							gas.NilBinds,
							gas.NilHooks,
							gas.NilHandlers,
							"article",
							map[string]string{
								"id": "article",
								"style": `border: 1px solid #dedede;padding: 2px 4px;margin-top:12px;`,
							})})
			})
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
