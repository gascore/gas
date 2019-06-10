package main

import (
	c "github.com/gascore/gas/examples/main/components"
	r "github.com/gascore/gas/std/router"
)

var ctx = &r.Ctx{

	Routes: []r.Route{
		{
			Name:      "home",
			Component: c.Home,
			Path:      "/",
			Exact:     true,
		},
		{
			Name:     "links-redirect",
			Exact:    true,
			Path:     "/link",
			Redirect: "/links",
		},
		{
			Name:     "link-redirect",
			Exact:    true,
			Path:     "/links/",
			Redirect: "/links",
		},
		{
			Name:      "link",
			Component: c.Link,
			Path:      "/links/:name",
		},
		{
			Name:           "todo-redirect",
			Exact:          true,
			Path:           "/todo",
			Redirect:       "todo-list",
			RedirectParams: map[string]string{"type": "all"},
		},
		{
			Name:     "todo-redirect",
			Exact:    true,
			Path:     "/todo/",
			Redirect: "/todo/all",
		},
		{
			Name:      "todo-list",
			Component: c.TodoList,
			Path:      "/todo/:type",
		},
		{
			Name:      "links",
			Component: c.Links,
			Path:      "/links",
			Exact:     true,
		},
		{
			Name:      "components",
			Component: c.Components,
			Path:      "/components",
		},
		{
			Name:      "about",
			Component: c.About,
			Path:      "/about",
		},
	},

	Settings: r.Settings{
		// BaseName: "/router",
		HashMode: true,
		GetUserConfirmation: func() bool {
			return false
		},
		ForceRefresh: !r.SupportHistory(),
	},
}
