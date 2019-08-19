package gas

import "testing"

func TestUpdateChildes(t *testing.T) {
	tree := NE(
		&E{Tag: "main", RC: GetEmptyRenderCore()},
		NE(
			&E{Tag: "div"},
			"wow",
		),
		NE(
			&E{Tag: "h1", Attrs: func() Map { return Map{"id": "wow"} }},
			"Title",
		),
		NE(
			&E{Tag: "p"},
			"Lorem ipsum dolore",
			" ",
			NE(
				&E{Tag: "i", Attrs: func() Map { return Map{"id": "lol", "class": "some"} }},
				"opsum",
			),
			NE(
				&E{Tag: "b"},
				"fote",
				NE(
					&E{Tag: "i"},
					"wate",
				),
				NE(
					&E{HTML: func() string { return "<h2>some</h2>" }},
				),
			),
		),
	)
	tree.UpdateChildes()

	isChildValid := func(i interface{}) {
		e, ok := i.(*E)
		if !ok {
			t.Errorf("invalid first child type: %T", tree.Childes[0])
		}
		if len(e.Tag) == 0 || len(e.UUID) == 0 || e.RC == nil || e.Parent == nil {
			t.Errorf("invalid child *Element: %v", e)
		}
	}

	isChildValid(tree.Childes[0])
	isChildValid(tree.Childes[1])
	isChildValid(tree.Childes[2])
	isChildValid(tree.Childes[2].(*E).Childes[2])
	isChildValid(tree.Childes[2].(*E).Childes[3])
	isChildValid(tree.Childes[2].(*E).Childes[3].(*E).Childes[1])
	isChildValid(tree.Childes[2].(*E).Childes[3].(*E).Childes[2])

	if tree.Childes[2].(*E).Childes[3].(*E).Childes[2].(*E).RHTML == "" {
		t.Error("invalid html directive")
	}
}

func TestUpdateElementChildes(t *testing.T) {
	var nodes []*RenderTask

	root := &exampleRoot{
		msg:     "no",
		counter: 5,
	}
	c := &C{
		RC: &RenderCore{
			BE: emptyBackEnd{
				logger: func(newNodes []*RenderTask) {
					nodes = append(nodes, newNodes...)
				},
			},
		},
		Root: root,
	}
	root.c = c

	f := func(i int) {
		if len(nodes) != i {
			t.Errorf("not enough render nodes, want: 1, but got: %d, nodes: %v", len(nodes), nodes)
		}
		nodes = []*RenderTask{}
	}

	el := root.c.Init()
	el.Update()
	f(4)

	root.counter = 6
	root.msg = "wow"
	root.c.Update()
	f(1)
}
