package gas

import "testing"

func TestNew(t *testing.T) {
	c := &C{Root: &exampleRoot{}}

	gas, err := New(GetEmptyBackend(), "app", c)
	if err != nil {
		t.Error("error in gas.New()")
	}

	gas.App.UpdateChildes()

	if len(gas.App.Childes) == 0 {
		t.Error("app has no childes")
	}

	child, ok := gas.App.Childes[0].(*E)
	if !ok {
		t.Error("invalid app first child type")
	}

	if child.Attrs["class"] != "wow" {
		t.Error("app has wrong start point")
	}
}

type exampleRoot struct{}

func (root *exampleRoot) Render() []interface{} {
	return CL(
		NE(
			&E{
				Attrs: map[string]string{
					"class": "wow",
				},
			},
			"text",
		),
		NE(
			&E{
				Tag: "h1",
			},
			"header",
		),
	)
}
