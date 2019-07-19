package gas

import "testing"

func TestChanged(t *testing.T) {
	f := func(a, b interface{}, isChanged bool, haveError bool) {
		ic, err := Changed(a, b)

		if err == nil {
			if haveError {
				t.Error("except error, but don't got it")
			}

			if ic != isChanged {
				t.Errorf("invalid Changed result except: %t, got: %t", isChanged, ic)
			}
		} else if !haveError {
			t.Errorf("error from Changed: %s", err.Error())
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
	f("1", "1", false, false)

	// various types
	f("1", 1, true, false)
	f(&C{}, &E{}, true, false)

	// not supported type
	f(Gas{}, Gas{}, true, true)

	// empty stucts
	f(&E{}, &E{}, false, false)
	f(&C{}, &C{}, false, false)

	// Component.Element
	f(c1, c2, true, false)
	f(c1, c1, false, false)

	// Element.Component
	f(&E{Component: c1}, &E{Component: c2}, true, false)
	f(&E{Component: c1}, &E{Component: c1}, false, false)
	f(&E{Component: c1}, &E{}, true, false)

	// attrs
	f(&E{Attrs: map[string]string{"id": "wow"}}, &E{Attrs: map[string]string{"id": "wow"}}, false, false)
	f(&E{Attrs: map[string]string{"id": "wow"}}, &E{Attrs: map[string]string{"id": "wow", "class": "lol"}}, true, false)
	f(&E{Attrs: map[string]string{"id": "wow"}}, &E{Attrs: map[string]string{"id": "lol"}}, true, false)

	// watchers
	f(&C{Watchers: map[string]Watcher{"wow": func(a interface{}, e Object) (string, error) { return "wow", nil }}}, &C{}, true, false)
	f(&C{Watchers: map[string]Watcher{"wow": func(a interface{}, e Object) (string, error) { return "wow", nil }}}, &C{Watchers: map[string]Watcher{"wow": func(a interface{}, e Object) (string, error) { return "wow", nil }}}, true, false)

	// hooks
	m1 := func() error { return nil }
	m2 := func() (bool, error) { return false, nil }

	f(&C{Hooks: Hooks{Created: m1}}, &C{Hooks: Hooks{Created: m1}}, false, false)
	f(&C{Hooks: Hooks{Created: m1}}, &C{Hooks: Hooks{}}, true, false)

	f(&C{Hooks: Hooks{BeforeCreated: m2}}, &C{Hooks: Hooks{BeforeCreated: m2}}, false, false) // with control
	f(&C{Hooks: Hooks{BeforeCreated: m2}}, &C{Hooks: Hooks{}}, true, false)                   // with control
}
