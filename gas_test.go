package gas

import (
	"testing"
)

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

type exampleRoot struct {
	c *C

	msg     string
	counter int
}

func (root *exampleRoot) Render() []interface{} {
	return CL(
		NE(
			&E{
				Attrs: map[string]string{
					"class": "wow",
				},
			},
			root.msg,
			func() interface{} {
				if root.msg == "wow" {
					return "Message is wow"
				}
				return nil
			},
		),
		NE(
			&E{
				Tag: "h1",
			},
			root.counter,
		),
		func() interface{} {
			if root.counter == 5 {
				return "Number is 5!"
			}
			return nil
		}(),
		NE(
			&E{},
			func() interface{} {
				if root.counter == 6 {
					return NE(&E{Tag: "i"}, "Counter is 6")
				}
				return NE(&E{Tag: "b"}, "Counter isn't 6")
			}(),
		),
	)
}
