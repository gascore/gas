package gas

import "testing"

func TestNew(t *testing.T) {
	gas, err := New(
		getEmptyBackend(),
		"app",
		&C{

		}, func(this *Component) []interface{} {
			return ToGetComponentList()
		})
	if err != nil {
		t.Error("in New function error can be thrown only by backend")
	}

	if gas.App.be != getEmptyBackend() {
		t.Error("invalid backend")
		return
	}

	if len(gas.App.Attrs) == 0 {
		t.Error("app has empty attributes")
		return
	}

	if len(gas.App.Attrs) == 0 {
		t.Error("app has empty attributes")
		return
	}

	if gas.App.Attrs["id"] != "app" {
		t.Error("app has wrong start point")
	}

	if gas.GetElement() == nil {
		t.Error("wtf?!")
	}
}
