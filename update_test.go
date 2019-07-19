package gas

import "testing"

func TestUpdateChildes(t *testing.T) {
	tree := NE(
		&E{Tag:"main",RC:GetEmptyRenderCore()},
		NE(
			&E{Tag:"div"},
			"wow",
		),
		NE(
			&E{Tag:"h1",Binds:map[string]Bind{"id": func()string {return "wow"}}},
			"Title",
		),
		NE(
			&E{Tag:"p"},
			"Lorem ipsum dolore",
			" ",
			NE(
				&E{Tag: "i",Binds:map[string]Bind{"id": func()string {return "lol"},"class": func()string {return "some"}}},
				"opsum",
			),
			NE(
				&E{Tag: "b"},
				"fote",
				NE(
					&E{Tag:"i"},
					"wate",
				),
				NE(
					&E{HTML:HTMLDirective{Render:func()string{return "<h2>some</h2>"}}},
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
		if len(e.Binds) != 0 {
			if len(e.RenderedBinds) == 0 {
				t.Errorf("invalid child RBinds: %v", e)
			}

			for key, val := range e.Binds {
				if val() != e.RenderedBinds[key] {
					t.Errorf("invalid child RBinds value key: %v, val: %v", key, e.RenderedBinds[key])
				}
			} 
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
