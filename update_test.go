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
			&E{Tag: "h1", func() map[string]string {return map[string]string{"id": "wow"}}},
			"Title",
		),
		NE(
			&E{Tag: "p"},
			"Lorem ipsum dolore",
			" ",
			NE(
				&E{Tag: "i", Attrs: func() map[string]string {return map[string]string{"id": "lol", "class": "some"}},
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
					&E{HTML: HTMLDirective{Render: func() string { return "<h2>some</h2>" }}},
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

	if tree.Childes[2].(*E).Childes[3].(*E).Childes[2].(*E).HTML.Rendered == "" {
		t.Error("invalid html directive")
	}
}

func TestUpdateElementChildes(t *testing.T) {
	var nodes []*RenderNode

	root := &exampleRoot{
		msg:     "no",
		counter: 5,
	}
	c := &C{
		RC: &RenderCore{
			BE: emptyBackEnd{
				logger: func(node *RenderNode) {
					// fmt.Println("NODE", node)
					nodes = append(nodes, node)
				},
			},
		},
		Root: root,
	}
	root.c = c

	el := root.c.Init()
	el.Update()

	// fmt.Println(el.Childes[2].(*E).Childes)
	// fmt.Println(el.Childes[3].(*E).Childes)
	if len(nodes) != 4 {
		t.Errorf("not enough render nodes, want: 4, but got: %d, nodes: %v", len(nodes), nodes)
	}
	// fmt.Println(el.Childes[3].(*E).Childes)
	// fmt.Println(el.Childes)
	nodes = []*RenderNode{}

	root.counter = 6
	root.msg = "wow"

	root.c.Update()
	// fmt.Println(el.Childes[2].(*E).Childes[0])
	if len(nodes) != 1 {
		// fmt.Println(nodes[0])
		t.Errorf("not enough render nodes, want: 1, but got: %d, nodes: %v", len(nodes), nodes)
	}
	nodes = []*RenderNode{}

	// fmt.Println(nodes)
	// fmt.Println(nodes[0])
	// fmt.Println(nodes[1])

	// 	fmt.Println(el.Childes[2].(*E).Childes)
	// 	fmt.Println(el.Childes[3].(*E).Childes)
}
