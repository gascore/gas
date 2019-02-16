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
				"changeType": func(this *Component, values ...interface{}) error {
					if len(values) == 0 {
						return errors.New("method values are nil")
					}
					if len(values) != 1 {
						return errors.New("not one value in method values")
					}

					val, ok := values[0].(int)
					if !ok {
						return errors.New("invalid value type")
					}

					err := this.SetData("type", val)
					if err != nil {
						return err
					}

					return nil
				},
			},
		},
		func(this *Component) []interface{} {
			return ToGetComponentList(
				"wow",
				this.GetData("type"))
		})
	c.be = getEmptyBackend()

	err := c.Method("changeType", 1)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}

	if c.GetData("type").(int) != 1 {
		t.Errorf("method was not called")
		return
	}

	err = nil
	err = c.Method("invalidMethodName", 1, 2, 3)
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
					val, ok := this.GetData("type").(int)
					this.WarnIfNot(ok)

					return val + 2, nil
				},
			},
		},
		func(this *Component) []interface{} {
			return ToGetComponentList(
				"wow",
				this.GetData("type"))
		})
	c.be = getEmptyBackend()

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
