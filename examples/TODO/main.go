package main

import (
	"errors"
	"fmt"
	"github.com/Sinicablyat/dom"
	"github.com/Sinicablyat/gas"
)

// Example application #11
//
// 'todo' shows how you how to build basic TODO example
func main() {
	app, err :=
		gas.New(
			"app",
			func(p *gas.Component) interface{} {
				return gas.NewComponent(
					&gas.Component{
						ParentC: p,
						Data: map[string]interface{}{
							"currentList": "0",
							"currentText": "",

							"current": []interface{}{},
							"done":    []interface{}{},
							"deleted": []interface{}{},
						},
						Methods: map[string]gas.Method{
							"delete": func(this *gas.Component, values ...interface{}) error {
								i, ok := values[0].(int)
								if !ok {
									return errors.New("invalid index")
								}

								appendToDeleted, ok := values[1].(bool)
								if !ok {
									return errors.New("invalid appendToDeleted")
								}

								list, ok := this.GetData("current").([]interface{})
								if !ok {
									return errors.New("invalid current list")
								}
								removedItem := list[i]

								err := this.DoWithUpdate(func() error {return this.SetData("current", remove(list, i))})
								if err != nil {
									return err
								}

								err = this.SetData("currentText", "")
								if err != nil {
									return err
								}

								if appendToDeleted {
									err = this.Method("append", "deleted", removedItem)
									if err != nil {
										return err
									}
								}

								return nil
							},
							"append": func(this *gas.Component, values ...interface{}) error {
								listTypeS, ok := values[0].(string)
								if !ok {
									return errors.New("invalid list type")
								}

								newTask, ok := values[1].(string)
								if !ok {
									return errors.New("invalid task")
								}

								list := this.GetData(listTypeS).([]interface{})
								list  = append(list, newTask)

								err := this.SetData(listTypeS, list)
								if err != nil {
									return err
								}

								if listTypeS == "current" {
									gas.WarnError(this.SetData("currentText", ""))
								}

								return nil
							},
							"markAsDone": func(this *gas.Component, values ...interface{}) error {
								i, ok := values[0].(int)
								if !ok {
									return errors.New("invalid index")
								}

								list := this.GetData("current").([]interface{})

								item := list[i]

								err := this.Method("append", "done", item)
								if err != nil {
									return err
								}

								err = this.Method("delete", i, false)
								if err != nil {
									return err
								}

								return nil
							},
							"edit": func(this *gas.Component, values ...interface{}) error {
								i, ok := values[0].(int)
								if !ok {
									return errors.New("invalid index")
								}

								newValue, ok := values[1].(string)
								if !ok {
									return errors.New("invalid new value")
								}

								list, ok := this.GetData("current").([]interface{})
								if !ok {
									return errors.New("invalid current list")
								}

								err := this.SetData("current", addItem(list, i, newValue))
								if err != nil {
									return err
								}

								err = this.Rerender()
								if err != nil {
									return err
								}

								return nil
							},
						},
						Tag: "div",
						Attrs: map[string]string{
							"id": "todo",
						},
					},
					func(this *gas.Component) interface{} {
						return gas.NewComponent(
							&gas.Component{
								ParentC:this,
								Tag:"style",
								Attrs: map[string]string{"type": "text/css"},
								Directives: gas.Directives{
									HTML: gas.HTMLDirective{
										Render: func(this2 *gas.Component) string {
											return `
#todo {
	width: 50%;
	margin: 0 auto;
}

#main {
	border: 1px solid #dedede;
	padding: 8px 16px 8px 16px;
}

ul {
	padding: 0;
	list-style-type: none;
	padding-left: 0;
	margin-left: 0;
}

ul li {
	display: flex;
	margin: 4px 0 8px 0;
	padding: 4px 8px;
	border-bottom: 1px solid #dedede;

	font-size: 18px;
}

ul li button {
	border: 0;
	padding: 0;
	background-color: inherit;
	cursor: pointer;
}
ul li button#submit:hover, button#submit:focus {
	color: #009966;
}
ul li button#delete:hover, button#delete:focus {
	color: #ff0033;
}

ul li button#submit {
	margin: 0 12px 0 0;
}

ul li button#delete {
	margin: 0 0 0 auto;
}

nav {
	margin-bottom: 8px;
}

nav button {
	margin-right: 6px;
	border: 0;
	padding: 0;
	color: #009966;
	background-color: inherit;
	text-decoration: underline;
	cursor: pointer;
}
nav button:focus, nav button:hover, nav button.active {
	color: #00CC99;
}

#new {
	width: auto;
}

footer {
	margin-top: 18px;
	color: gray;
	font-size: 12px;
	text-align: center;
}
footer div {
	margin-bottom: 4px;
}
footer a {
	margin: 0 4px;
	color: inherit;
}

`
						}}}})
					},
					func (this *gas.Component) interface{} {
						return gas.NewComponent(
							&gas.Component{
								ParentC:this,
								Tag: "div",
								Attrs: map[string]string{
									"id": "main",
								},
							},
							func (p *gas.Component) interface{} {
								return gas.NewComponent(
									&gas.Component{
										ParentC: p,
										Tag: "nav",
									},
									func(p *gas.Component) interface{} {
										return gas.NewComponent(
											&gas.Component{
												ParentC: p,
												Tag: "button",
												Handlers: map[string]gas.Handler{
													"click": func(p *gas.Component, e dom.Event) {
														gas.WarnError(this.SetData("currentList", "0"))
													},
												},
												Binds: map[string]gas.Bind{
													"class": func(p *gas.Component) string {
														if this.GetData("currentList").(string) == "0" {
															return "active"
														}
														return ""
													},
												},
											},
											func(p *gas.Component) interface{} {
												return "Current"
											})
									},
									func(p *gas.Component) interface{} {
										return gas.NewComponent(
											&gas.Component{
												ParentC: p,
												Tag: "button",
												Handlers: map[string]gas.Handler{
													"click": func(p *gas.Component, e dom.Event) {
														gas.WarnError(this.SetData("currentList", "1"))
													},
												},
												Binds: map[string]gas.Bind{
													"class": func(p *gas.Component) string {
														if this.GetData("currentList").(string) == "1" {
															return "active"
														}
														return ""
													},
												},
											},
											func(p *gas.Component) interface{} {
												return "Completed"
											})
									},
									func(p *gas.Component) interface{} {
										return gas.NewComponent(
											&gas.Component{
												ParentC: p,
												Tag: "button",
												Handlers: map[string]gas.Handler{
													"click": func(p *gas.Component, e dom.Event) {
														gas.WarnError(this.SetData("currentList", "2"))
													},
												},
												Binds: map[string]gas.Bind{
													"class": func(p *gas.Component) string {
														if this.GetData("currentList").(string) == "2" {
															return "active"
														}
														return ""
													},
												},
											},
											func(p *gas.Component) interface{} {
												return "Deleted"
											})
									},)
							},
							func (p *gas.Component) interface{} {
								return gas.NewComponent(
									&gas.Component{
										ParentC: p,
										Directives: gas.Directives{
											Model: gas.ModelDirective{
												Data: "currentText",
												Component: this,
											},
										},
										Tag: "input",
										Handlers: map[string]gas.Handler{
											"keyup.enter": func(p *gas.Component, e dom.Event) {
												currentText := this.GetData("currentText").(string)
												if len(currentText) == 0 {
													return
												}

												gas.WarnError(this.Method("append", "current", currentText))
											},
										},
										Attrs: map[string]string{
											"id": "new",
											"placeholder": "New task",
										},
									})
							},
							func (p *gas.Component) interface{} {
								return gas.NewComponent(
									&gas.Component{
										ParentC: p,
									},
									// Because i don't need wrap `this` and ul `this` i can overwrite this variable
									func(p *gas.Component) interface{} {
										return gas.NewComponent(
											&gas.Component{
												ParentC: p,
												Directives:gas.Directives{
													Show: func(p *gas.Component) bool {
														return this.GetData("currentList") == "0"
													},
												},
												Tag: "ul",
												Attrs: map[string]string{
													"id": "list__current",
													"class": "list",
												},
											},
											func(p *gas.Component) interface{} {
												return getLi(p, this, 0)
											})
									},
									func(p *gas.Component) interface{} {
										return gas.NewComponent(
											&gas.Component{
												ParentC: p,
												Tag: "ul",
												Directives:gas.Directives{
													Show: func(p *gas.Component) bool {
														return this.GetData("currentList") == "1"
													},
												},
												Attrs: map[string]string{
													"id": "list__done",
													"class": "list",
												},
											},
											func(p *gas.Component) interface{} {
												return getLi(p, this, 1)
											})
									},
									func(p *gas.Component) interface{} {
										return gas.NewComponent(
											&gas.Component{
												ParentC: p,
												Tag: "ul",
												Directives:gas.Directives{
													Show: func(p *gas.Component) bool {
														return this.GetData("currentList") == "2"
													},
												},
												Attrs: map[string]string{
													"id": "list__deleted",
													"class": "list",
												},
											},
											func(p *gas.Component) interface{} {
												return getLi(p, this, 2)
											})
									})
							},)
					},
					func (this *gas.Component) interface{} {
						return gas.NewComponent(
							&gas.Component{
								ParentC:this,
								Tag:"footer",
							},
							func(p *gas.Component) interface{} {
								return gas.NewComponent(
									&gas.Component{
										ParentC:p,
										Tag:"div",
									},
									func(p *gas.Component) interface{} {
										return "Double-click to edit a task"
									})
							},
							func(p *gas.Component) interface{} {
								return gas.NewComponent(
									&gas.Component{
										ParentC:p,
										Tag:"div",
									},
									func(p *gas.Component) interface{} {
										return "Created by"
									},
									func(p *gas.Component) interface{} {
										return gas.NewComponent(
											&gas.Component{
												Tag: "a",
												Attrs: map[string]string{
													"href": "https://sinicablyat.github.io/",
													"target": "_blank",
												},
											},
											func(p *gas.Component) interface{} {
												return "Noskov Artem"
											})
									},
									func(p *gas.Component) interface{} {
										return "with"
									},
									func(p *gas.Component) interface{} {
										return gas.NewComponent(
											&gas.Component{
												Tag: "a",
												Attrs: map[string]string{
													"href": "https://sinicablyat.github.io/gas",
													"target": "_blank",
												},
											},
											func(p *gas.Component) interface{} {
												return "GAS"
											})
									},
									func(p *gas.Component) interface{} {
										return "and love"
									},)
							})
					})
			},
		)
	must(err)

	err = app.Init()
	must(err)
	gas.KeepAlive()
}

func getLi(p *gas.Component, this *gas.Component, listType int) interface{} {
	// listType: 0 - current, 1 - done, 2 - deleted tasks
	var listTypeS string
	switch listType {
	case 0:
		listTypeS = "current"
	case 1:
		listTypeS = "done"
	case 2:
		listTypeS = "deleted"
	}

	return gas.NewComponent(
		&gas.Component{
			ParentC: p,
			Tag: "li",
			Data: map[string]interface{}{
				"isEditing": false,
				"newValue": "no",
			},
			Directives:gas.Directives{
				For:gas.ForDirective{
					Data: listTypeS,
					Component: this,
					Render: func(i int, value interface{}, p *gas.Component) []gas.GetComponent {
						return gas.ToGetComponentList(
							func(p *gas.Component) interface{} {
								return gas.NewComponent(
									&gas.Component{
										Tag: "button",
										Directives:gas.Directives{
											If: func(p *gas.Component) bool {
												return listType == 0
											},
										},
										Handlers: map[string]gas.Handler{
											"click": func(this5 *gas.Component, e dom.Event) {
												gas.WarnError(this.Method("markAsDone", i))
											},
										},
										Attrs: map[string]string{
											"id": "submit",
										},
									},
									func(this3 *gas.Component) interface{} {
										return gas.NewComponent(&gas.Component{Tag: "i", Attrs: map[string]string{"class": "fas fa-check"}})
									})
							},
							func(this2 *gas.Component) interface{} {
								return gas.NewComponent(
									&gas.Component{
										Tag: "i",
										Directives:gas.Directives{
											If: func(p *gas.Component) bool {
												return !this2.GetData("isEditing").(bool)
											},
										},
										Handlers: map[string]gas.Handler{
											"dblclick": func(p *gas.Component, e dom.Event) {
												gas.WarnError(this2.SetData("newValue", value))
												gas.WarnError(this2.SetData("isEditing", true))
											},
										},
									},
									func(p *gas.Component) interface{} {
										return fmt.Sprintf("%s", value)
									})
							},
							func(this2 *gas.Component) interface{} {
								return gas.NewComponent(
									&gas.Component{
										Tag: "input",
										Attrs: map[string]string{
											"style": "margin-right: 8px",
										},
										Directives:gas.Directives{
											If: func(p *gas.Component) bool {
												return this2.GetData("isEditing").(bool)
											},
											Model:gas.ModelDirective{
												Component: this2,
												Data: "newValue",
											},
										},
										Handlers: map[string]gas.Handler{
											"keyup.enter": func(p *gas.Component, e dom.Event) {
												gas.WarnError(this2.SetData("isEditing", false))
												gas.WarnError(this.Method("edit", i, this2.GetData("newValue")))
											},
										},
									},
									func(p *gas.Component) interface{} {
										return fmt.Sprintf("%s", value)
									})
							},
							func(p *gas.Component) interface{} {
								return gas.NewComponent(
									&gas.Component{
										Tag: "button",
										Directives:gas.Directives{
											If: func(p *gas.Component) bool {
												return listType == 0
											},
										},
										Handlers: map[string]gas.Handler{
											"click": func(this5 *gas.Component, e dom.Event) {
												gas.WarnError(this.Method("delete", i, true))
											},
										},
										Attrs: map[string]string{
											"id": "delete",
										},
									},
									func(this3 *gas.Component) interface{} {
										return gas.NewComponent(&gas.Component{Tag: "i", Attrs: map[string]string{"class": "fas fa-trash-alt "}})
									})
							},)
					},
				},
			},
		})
}

func remove(a []interface{}, i int) []interface{} {
	copy(a[i:], a[i+1:]) // Shift a[i+1:] left one index
	a[len(a)-1] = ""     // Erase last element (write zero value)
	a = a[:len(a)-1]     // Truncate slice

	return a
}

func addItem(a []interface{}, i int, newValue string) []interface{} {
	a[i] = newValue
	return a
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
