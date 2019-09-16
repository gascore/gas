package gas

import (
	"testing"
)

func TestNew(t *testing.T) {
	c := &C{Root: &exampleRoot{}}

	gas := New(c, GetEmptyBackend())

	gas.Component.Update()

	el := gas.Component.Element
	if el == nil {
		t.Error("component element is nil")
	}

	if len(el.Childes) == 0 {
		t.Error("app has no childes")
	}

	child, ok := el.Childes[0].(*E)
	if !ok {
		t.Error("invalid app first child type")
	}

	if child.Attrs()["class"] != "wow" {
		t.Error("app has wrong start point")
	}
}

type exampleRoot struct {
	c *C

	msg     string
	counter int
}

func (root *exampleRoot) Render() *Element {
	return NE(
		&E{},
		NE(
			&E{
				Attrs: func() Map {
					return Map{
						"class": "wow",
					}
				},
			},
			root.msg,
			func() interface{} {
				if root.msg == "wow" {
					return "Message is wow"
				}
				return nil
			}(),
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
