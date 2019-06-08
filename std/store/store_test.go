package store

import (
	"github.com/gascore/gas"
	"github.com/pkg/errors"
	"testing"
)

func TestNew(t *testing.T) {
	var (
		middlewareWasCalled bool
		onCreateWasCalled   bool
		beforeEmitWasCalled bool
		afterEmitWasCalled  bool
	)

	bStore := &Store{
		Data: map[string]interface{}{
			"counter": 0,
		},

		Handlers: map[string]Handler{
			"addToCounter": func(s *Store, values ...interface{}) (updateMap map[string]interface{}, err error) {
				if len(values) == 0 {
					return nil, errors.New("invalid values")
				}

				value, ok := values[0].(int)
				if !ok {
					return nil, errors.New("invalid value type")
				}

				counter := s.Get("counter").(int)
				return map[string]interface{}{
					"counter": counter+value,
				},nil
			},
		},

		MiddleWares: []MiddleWare{
			{
				Prefix: "addTo",
				Hook: func(s *Store, values []interface{}) error {
					middlewareWasCalled = true
					return nil
				},
			},
		},

		OnCreate:[]OnCreateHook{
			func(s *Store) error {
				onCreateWasCalled = true
				return nil
			},
		},
		BeforeEmit: []BeforeEmitHook{
			func(s *Store, eventName string, values []interface{}) error {
				beforeEmitWasCalled = true
				return nil
			},
		},
		AfterEmit: []AfterEmitHook{
			func(s *Store, eventName string, updatesMap map[string]interface{}, values []interface{}) error {
				afterEmitWasCalled = true
				return nil
			},
		},
	}

	// New
	s, err := New(bStore)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}

	if !onCreateWasCalled {
		t.Error("onCreate was not called")
		return
	}

	// Add
	registeredComponent := s.RegisterComponent(gas.NC(
		&gas.C{
			Tag: "p",
			Attrs: map[string]string{
				"class": "component",
			},
			RC: gas.GetEmptyRenderCore(),
		},
		func(this *gas.Component) []interface{} {
			return []interface{}{
				s.Get("counter"),
			}
		},
	), nil)
	if registeredComponent == nil {
		t.Errorf("store RegisterComponent result is nil")
	}

	// Mounted
	err = gas.RunMountedIfCan(registeredComponent)
	if err != nil {
		t.Errorf("unexpected error in RunMountedIfCan: %s", err.Error())
		return
	}

	if len(s.subscribers) == 0 {
		t.Error("component was not added to store subscribers")
		return
	}

	// Emit
	err = s.Emit("addToCounter", 1)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}

	if !beforeEmitWasCalled {
		t.Error("beforeEmit was not called")
		return
	}

	if !afterEmitWasCalled {
		t.Error("afterEmit was not called")
		return
	}

	if !middlewareWasCalled {
		t.Error("middleware was not called")
		return
	}

	// BeforeDestroy
	err = gas.RunWillDestroyIfCan(registeredComponent)
	if err != nil {
		t.Errorf("unexpected error in RunWillDestroyIfCan: %s", err.Error())
		return
	}

	if len(s.subscribers) != 0 {
		t.Error("component was not removed from store subscribers")
		return
	}
}

func TestIsRoot(t *testing.T) {
	s, err := New(&Store{
		Data: map[string]interface{}{
			"counter": 0,
		},
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}

	c := gas.NC(
		&gas.C{
			Tag: "p",
			Attrs: map[string]string{
				"class": "component",
			},
		},
		func(this *gas.Component) []interface{} {
			return []interface{}{
				"component",
			}
		},
	)

	var emptyChildes gas.GetComponentChildes

	componentInStore := gas.NC(
		&gas.C{
			Tag: "div",
		}, emptyChildes)

	s.subscribers = append(s.subscribers, Sub{componentInStore, nil})

	data := []struct{
		parent *gas.Component
		isRoot bool
	}{
		{
			parent: gas.NC(
				&gas.C{
					Tag: "h1",
					Parent: gas.NC(
						&gas.C{
							Tag: "div",
						}, emptyChildes),
				}, emptyChildes,
			),
			isRoot: true,
		},
		{
			parent: gas.NC(
				&gas.C{
					Tag: "h1",
					Parent: gas.NC(
						&gas.C{
							Tag: "div",
							Parent: componentInStore,
						}, emptyChildes),
				}, emptyChildes,
			),
			isRoot: false,
		},
		{
			parent: gas.NE(
				&gas.C{
					Tag: "h1",
					Parent: gas.NC(
						&gas.C{
							Tag: "div",
							Parent: componentInStore,
						}, emptyChildes),
				},
			),
			isRoot: false,
		},
	}

	for _, el := range data {
		c.Parent = el.parent
		isRoot := s.isRoot(c)
		if isRoot != el.isRoot {
			t.Errorf("invalid isRoot result want: %t, got: %t", el.isRoot, isRoot)
		}
	}
}
