package gas

import (
	"testing"
)

func TestChanged(t *testing.T) {
	var counter int
	f := func(a, b interface{}, isChanged, canGoDeeper bool, haveError bool) {
		counter++
		ic, cgd, err := Changed(a, b)

		if err == nil {
			if haveError {
				t.Errorf("%d: except error, but don't got it", counter)
			}

			if ic != isChanged {
				t.Errorf("%d: invalid Changed isChanged result except: %t, got: %t", counter, isChanged, ic)
			}

			if cgd != canGoDeeper {
				t.Errorf("%d: invalid Changed canGoDeeper resul except: %t, got: %t", counter, canGoDeeper, cgd)
			}
		} else if !haveError {
			t.Errorf("%d: error from Changed: %s", counter, err.Error())
		}
	}

	e1 := &E{
		Tag: "h1",
	}

	e2 := &E{
		Tag: "h2",
	}

	c1 := &C{
		ElementIsImportant: true,
		Element:            e1,
	}

	c2 := &C{
		ElementIsImportant: true,
		Element:            e2,
	}

	// atomic types
	f("1", "1", false, false, false)
	f(1, 1, false, false, false)
	f(1.228, 1.228, false, false, false)

	// various types
	f("1", 1, true, false, false)
	f(&C{}, &E{}, true, false, false)

	// not supported type
	f(Gas{}, Gas{}, true, false, true)

	// empty stucts
	f(&E{}, &E{}, false, true, false)
	f(&C{}, &C{}, false, true, false)

	// Component.Element
	f(c1, c2, true, false, false)
	f(c1, c1, false, true, false)

	// Element.Component
	f(&E{Component: c1}, &E{Component: c2}, true, false, false)
	f(&E{Component: c1}, &E{Component: c1}, false, true, false)
	f(&E{Component: c1}, &E{}, true, false, false)

	eA1 := &E{Attrs: func() Map { return Map{"id": "wow"} }}
	eA1.RAttrs = eA1.Attrs()
	eA2 := &E{Attrs: func() Map { return Map{"id": "wow", "class": "lol"} }}
	eA2.RAttrs = eA2.Attrs()
	eA3 := &E{Attrs: func() Map { return Map{"id": "lol"} }}
	eA3.RAttrs = eA3.Attrs()

	// attrs
	f(eA1, eA1, false, true, false)
	f(eA1, eA2, true, true, false)
	f(eA1, eA3, true, true, false)

	// hooks
	m1 := func() error { return nil }
	m2 := func() (bool, error) { return false, nil }

	f(&C{Hooks: Hooks{Created: m1}}, &C{Hooks: Hooks{Created: m1}}, false, true, false)
	f(&C{Hooks: Hooks{Created: m1}}, &C{Hooks: Hooks{}}, true, false, false)

	f(&C{Hooks: Hooks{BeforeCreated: m2}}, &C{Hooks: Hooks{BeforeCreated: m2}}, false, true, false) // with control
	f(&C{Hooks: Hooks{BeforeCreated: m2}}, &C{Hooks: Hooks{}}, true, false, false)                  // with control
}

func TestDiffAttrs(t *testing.T) {
	f := func(newA, oldA Map, excp Map) {
		got := DiffAttrs(newA, oldA)
		if !compareAttributes(got, excp) {
			t.Errorf("error in DiffAttrs: excepted: %v, got: %v", excp, got)
		}
	}

	f(Map{"1": "1"}, Map{"1": "1"}, Map{})
	f(Map{"1": "1"}, Map{"1": "2"}, Map{"1":"1"})
	f(Map{"1": "3"}, Map{"1": "3"}, Map{})
	f(Map{"1": "1", "2": "2"}, Map{"1": "1"}, Map{"2": "2"})
	f(Map{"1": "1", "2": "2"}, Map{"1": "1", "2": "3"}, Map{"2": "2"})
	f(Map{"1": "1"}, Map{"1": "1", "2": "2"}, Map{"2": ""})
	f(Map{"1": "1", "2": "3"}, Map{"1": "1", "2": "2", "4": "4"}, Map{"2": "3", "4": ""})
	f(Map{"1": "1", "4": "4"}, Map{"1": "1", "2": "2", "3": "3"}, Map{"4": "4", "2": "", "3": ""})
}
