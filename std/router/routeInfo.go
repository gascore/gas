package router

import (
	"fmt"
	"strings"

	"github.com/gascore/dom"
	"github.com/gascore/dom/js"
	"github.com/gascore/gas"
)

type RouteInfo struct {
	Name string
	URL  string

	Params      map[string]string // /links/:foo => {"foo": "bar"}
	QueryParams map[string]string // /links?foo=bar => {"foo": "bar"}

	Route Route

	Ctx *Ctx
}

func (i RouteInfo) Push(path string, replace bool) {
	i.Ctx.Push(path, replace)
}
func (ctx *Ctx) Push(path string, replace bool) {
	if ctx.Settings.GetUserConfirmation != nil && ctx.Settings.GetUserConfirmation() {
		return
	}

	ctx.ChangeRoute(path, replace)

	dom.GetWindow().DispatchEvent(js.New("Event", changeRouteEvent))
}

func (i RouteInfo) PushDynamic(name string, params, queries map[string]string, replace bool) {
	i.Ctx.PushDynamic(name, params, queries, replace)
}
func (ctx *Ctx) PushDynamic(name string, params, queries map[string]string, replace bool) {
	ctx.Push(ctx.fillPath(name, params, queries), replace)
}

func (ctx *Ctx) fillPath(name string, params, queries map[string]string) string {
	route := ctx.getRoute(name)
	if route.Name == "" {
		return ""
	}

	path := route.Path

	for x := 0; x < 64; x++ {
		p1, name, p2 := splitPath(path)
		if len(name) == 0 {
			var queriesString string
			if queries != nil {
				queriesString = "?"
				for key, value := range queries {
					queriesString = queriesString + key + "=" + value + "&"
				}
				queriesString = strings.TrimSuffix(queriesString, "&") // remove last "&"
			}

			return path + queriesString
		}

		path = fmt.Sprintf("%s%s%s", p1, params[name], p2)
	}

	ctx.This.WarnError(fmt.Errorf("invalid path"))
	return path
}

func (ctx *Ctx) getRoute(name string) Route {
	for _, r := range ctx.Routes {
		if r.Name == name {
			return r
		}
	}

	ctx.This.WarnError(fmt.Errorf("undefined route: %s", name))
	return Route{}
}

func (ctx Ctx) link(getPath func() string, push func(*gas.Component, gas.Object), e gas.External) *gas.Component {
	return gas.NE(
		&gas.Component{
			Tag: "a",
			Attrs: map[string]string{
				"href": "#",
			},
			Binds: map[string]gas.Bind{
				"href": func() string {
					return getPath()
				},
			},
			Handlers: map[string]gas.Handler{
				"click":    beforePush(push),
				"keyup.13": beforePush(push),
				"keyup.32": beforePush(push),
			},
		},
		e.Body...)
}
func beforePush(push func(*gas.Component, gas.Object)) func(*gas.Component, gas.Object) {
	return func(this *gas.Component, event gas.Object) {
		push(this, event)
		event.Call("preventDefault")
	}
}

func (i RouteInfo) LinkStatic(to string, replace bool, e gas.External) *gas.Component {
	return i.Ctx.LinkStatic(to, replace, e)
}
func (ctx *Ctx) LinkStatic(to string, replace bool, e gas.External) *gas.Component {
	return ctx.link(
		func() string {
			return ctx.Settings.BaseName + to
		},
		func(this *gas.Component, e gas.Object) {
			ctx.Push(to, replace)
		},
		e)
}

func (i RouteInfo) LinkDynamic(name string, params, queries map[string]string, replace bool, e gas.External) *gas.Component {
	return i.Ctx.LinkDynamic(name, params, queries, replace, e)
}
func (ctx *Ctx) LinkDynamic(name string, params, queries map[string]string, replace bool, e gas.External) *gas.Component {
	return ctx.link(
		func() string {
			return ctx.Settings.BaseName + ctx.fillPath(name, params, queries)
		},
		func(this *gas.Component, e gas.Object) {
			ctx.PushDynamic(name, params, queries, replace)
		},
		e)
}
