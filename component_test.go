package gas

import (
	"reflect"
	"testing"
)

func TestComponentInit(t *testing.T) {
	f := func(c *C) {
		c.Root = &exampleRoot{}

		el := c.Init()

		if el.UUID == "" {
			t.Errorf("empty element UUID")
		}

		if el.Tag == "" {
			t.Errorf("empty element tag")
		}

		el.UpdateChildes()
		if len(el.Childes) == 0 {
			t.Error("element childes are null")
		}
	}

	f(&C{Element: &E{Tag: "h1", Attrs: func() Map { return Map{"id": "foo"} }}})
	f(&C{Element: &E{Attrs: func() Map { return Map{"id": "foo"} }}})
	f(&C{Element: &E{UUID: "custom_id"}})
}

func TestUnSpliceBody(t *testing.T) {
	f := func(in, out []interface{}) {
		formatedIn := UnSpliceBody(in)
		if !reflect.DeepEqual(formatedIn, out) {
			t.Errorf("formatedIn and out are not the same: want - %v, but got - %v", out, formatedIn)
		}
	}

	f(CL(1, 2, 3),
		CL(1, 2, 3))

	f(CL(1, nil, 3),
		CL(1, 3))

	f(CL(1, CL(1, 2), 3),
		CL(1, 1, 2, 3))

	f(CL(1, 2, 3, []*E{&E{Tag: "h1"}, &E{Tag: "h2"}}),
		CL(1, 2, 3, &E{Tag: "h1"}, &E{Tag: "h2"}))

	f(CL(1, 2, 3, []*C{&C{Root: &exampleRoot{}}}),
		CL(1, 2, 3, &C{Root: &exampleRoot{}}))
}

func TestRemoveStrings(t *testing.T) {
	f := func(in, out []interface{}) {
		formatedIn := RemoveStrings(in)
		if !reflect.DeepEqual(formatedIn, out) {
			t.Errorf("formatedIn and out are not the same: want - %v, but got - %v", out, formatedIn)
		}
	}

	f(CL(1, 2, "3"),
		CL(1, 2))

	f(CL(1, "2", "3"),
		CL(1))

	f(CL(1, 2, 3, "4", 5),
		CL(1, 2, 3, 5))

	f(CL(1, 2, 3, 4, 5),
		CL(1, 2, 3, 4, 5))
}

func TestEmptyRoot(t *testing.T) {
	e := &E{Tag: "h1"}
	root := &EmptyRoot{Element: e}
	if !reflect.DeepEqual(root.Render(), CL(e)) {
		t.Errorf("invalid EmptyRoot.Render() result")
	}
}

func TestLittleUtils(t *testing.T) {
	// IsComponent
	if !IsComponent(&C{}) {
		t.Error("invalid IsComponent result #1")
	}
	if IsComponent(&E{}) {
		t.Error("invalid IsComponent result #2")
	}

	// IsElement
	if !IsElement(&E{}) {
		t.Error("invalid IsElement result #1")
	}
	if IsElement(&C{}) {
		t.Error("invalid IsElement result #2")
	}

	var cI interface{} = &C{}
	if !IsComponent(I2C(cI)) {
		t.Error("invalid I2C result")
	}

	var eI interface{} = &E{}
	if !IsElement(I2E(eI)) {
		t.Error("invalid I2E result")
	}
}

func TestElementParentComponent(t *testing.T) {
	f := func(e *E, haveParent bool) {
		p := e.ParentComponent()
		if (p == nil || p.Component == nil) && haveParent {
			t.Errorf("element parent is nil: %v", e)
		}
		if p != nil && p.Component != nil && !haveParent {
			t.Errorf("element parent isn't nil: %v", e)
		}
	}

	e1 := NE(&E{})
	e2 := NE(&E{Component: &C{}}, e1)
	p := NE(&E{}, e2)
	p.UpdateChildes()

	f(e1, true)
	f(e2, false)
}
