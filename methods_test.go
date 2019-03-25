package gas

import (
	"errors"
	"testing"
)

func TestComponent_Method(t *testing.T) {
	c := NC(
		&C{
			Data: map[string]interface{}{
				"type": 0,
			},
			Methods: map[string]Method{
				"changeType": func(this *Component, values ...interface{}) (interface{}, error) {
					if len(values) == 0 {
						return nil, errors.New("method values are nil")
					}
					if len(values) != 1 {
						return nil, errors.New("not one value in method values")
					}

					val, ok := values[0].(int)
					if !ok {
						return nil, errors.New("invalid value type")
					}

					this.Data["type"] = val

					return nil, nil
				},
			},
		},
		func(this *Component) []interface{} {
			return ToGetComponentList(
				"wow",
				this.Get("type"))
		})
	c.RC = GetEmptyRenderCore()

	nil, err := c.MethodSafely("changeType", 1)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}

	if c.Get("type").(int) != 1 {
		t.Errorf("method was not called")
		return
	}

	_, err = c.MethodSafely("invalidMethodName", 1, 2, 3)
	if err == nil {
		t.Error("no error after calling nil method")
		return
	}
}

func TestComponent_Computed(t *testing.T) {
	c := NC(
		&C{
			Data: map[string]interface{}{
				"type": 0,
			},
			Computeds: map[string]Computed{
				"getTypePlus2": func(this *Component, values ...interface{}) (interface{}, error) {
					val, ok := this.Get("type").(int)
					this.WarnIfNot(ok)

					return val + 2, nil
				},
			},
		},
		func(this *Component) []interface{} {
			return ToGetComponentList(
				"wow",
				this.Get("type"))
		})
	c.RC = GetEmptyRenderCore()

	val := c.Computed("getTypePlus2")
	if val == nil {
		t.Error("invalid computed response")
		return
	}

	val = c.Computed("invalidMethodName", 1, 2, 3)
	if val != nil {
		t.Error("value not nil after calling nil method")
		return
	}
}
