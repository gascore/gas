package gas

import (
	"testing"
)

func TestRunMountedIfCan(t *testing.T) {
	c := NC(
		&C{
			Tag: "h1",
			Data: map[string]interface{}{
				"counter":      0,
				"childCounter": 0,
			},
			Hooks: Hooks{
				Mounted: func(this *Component) error {
					return this.SetData("counter", 1)
				},
			},
		},
		func(this *Component) []interface{} {
			return ToGetComponentList(
				NC(
					&C{
						Hooks: Hooks{
							Mounted: func(this2 *Component) error {
								return this.SetData("childCounter", 1)
							},
						},
					},
					func(this2 *Component) []interface{} {
						return []interface{}{}
					},
				),
			)
		},
	)
	c.BE = GetEmptyBackend()

	c.RChildes = c.Childes(c) // Run{HookName}IfCan works with rendered childes

	err := RunMountedIfCan(c)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if c.GetData("counter").(int) != 1 {
		t.Errorf("mounted not called in parent component")
	}

	if c.GetData("childCounter").(int) != 1 {
		t.Errorf("mounted not called in child component")
	}
}

func TestRunWillDestroyIfCan(t *testing.T) {
	c := NC(
		&C{
			Tag: "h1",
			Data: map[string]interface{}{
				"counter":      0,
				"childCounter": 0,
			},
			Hooks: Hooks{
				WillDestroy: func(this *Component) error {
					return this.SetData("counter", 1)
				},
			},
		},
		func(this *Component) []interface{} {
			return ToGetComponentList(
				NC(
					&C{
						Hooks: Hooks{
							WillDestroy: func(this2 *Component) error {
								return this.SetData("childCounter", 1)
							},
						},
					},
					func(this2 *Component) []interface{} {
						return []interface{}{}
					},
				),
			)
		},
	)
	c.BE = GetEmptyBackend()

	c.RChildes = c.Childes(c) // Run{HookName}IfCan works with rendered childes

	err := RunWillDestroyIfCan(c)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if c.GetData("counter").(int) != 1 {
		t.Errorf("willDestroy not called in parent component")
	}

	if c.GetData("childCounter").(int) != 1 {
		t.Errorf("willDestroy not called in child component")
	}
}

func TestRunUpdatedIfCan(t *testing.T) {
	c2 := NC(
		&C{Tag: "i"},
		func(this2 *Component) []interface{} {
			return []interface{}{}
		})

	c := NC(
		&C{
			Tag: "h1",
			Data: map[string]interface{}{
				"counter":      0,
				"childCounter": 0,
			},
			Hooks: Hooks{
				Updated: func(this *Component) error {
					return this.SetData("counter", 1)
				},
			},
		},
		func(this *Component) []interface{} {
			return ToGetComponentList(c2)
		},
	)
	c.BE = GetEmptyBackend()

	c.RChildes = c.Childes(c) // Run{HookName}IfCan works with rendered childes

	err := RunUpdatedIfCan(c2)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if c.GetData("counter").(int) != 1 {
		t.Errorf("updated not called in parent component")
	}
}

func TestRunBeforeUpdateIfCan(t *testing.T) {
	c2 := NC(
		&C{Tag: "i"},
		func(this2 *Component) []interface{} {
			return []interface{}{}
		})

	c := NC(
		&C{
			Tag: "h1",
			Data: map[string]interface{}{
				"counter":      0,
				"childCounter": 0,
			},
			Hooks: Hooks{
				BeforeUpdate: func(this *Component) error {
					return this.SetData("counter", 1)
				},
			},
		},
		func(this *Component) []interface{} {
			return ToGetComponentList(c2)
		},
	)
	c.BE = GetEmptyBackend()

	c.RChildes = c.Childes(c) // Run{HookName}IfCan works with rendered childes

	err := RunBeforeUpdateIfCan(c2)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if c.GetData("counter").(int) != 1 {
		t.Errorf("updated not called in parent component")
	}
}
