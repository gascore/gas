package gas

import (
	"github.com/pkg/errors"
	"testing"
)

func TestGetData(t *testing.T) {
	data := []struct {
		c     *Component
		query string
		isNil bool
	}{
		{
			c: &C{
				Data: map[string]interface{}{
					"hello_world": "Hello world!!",
				},
				be: getEmptyBackend(),
			},
			query: "hello_world",
			isNil: false,
		},
		{
			c: &C{
				Data: map[string]interface{}{},
				be:   getEmptyBackend(),
			},
			query: "hello_world",
			isNil: true,
		},
		{
			c:     &C{be: getEmptyBackend()}, // check will GetData throw panic with nil Data map
			query: "hello_world",
			isNil: true,
		},
	}

	for _, el := range data {
		data := el.c.GetData(el.query)
		if data == nil && !el.isNil {
			t.Errorf("data is nil want: %t, got: %t", el.isNil, data == nil)
		}
	}
}

func TestSetData(t *testing.T) {
	errInByeWatcher := errors.New("why do we have to say goodbye?")

	c1 := NC(&C{
		Data: map[string]interface{}{
			"hello_world": "Hello world!!",
			"bye_world":   "no",
		},
		Watchers: map[string]Watcher{
			"hello_world": func(this *Component, new interface{}, old interface{}) error {
				// do something here
				return nil
			},
			"bye_world": func(this *Component, new interface{}, old interface{}) error {
				return errInByeWatcher
			},
		},
		be: getEmptyBackend(),
	}, func(this *Component) []interface{} {
		return []interface{}{}
	})
	c1.be = getEmptyBackend()

	c2 := NC(&C{
		Data: map[string]interface{}{
			"type": 0,
		},
		Directives: Directives{
			HTML: HTMLDirective{
				Render: func(this *Component) string {
					var val string
					if this.GetData("type").(int) == 0 {
						val = `<a href="https://gascore.github.io/" target="_blank">project site</a>`
					} else {
						val = `<a href="https://noartem.github.io/" target="_blank">my site</a>`
					}

					return val
				},
			},
		},
	}, func(this *Component) []interface{} {
		return []interface{}{}
	})
	c2.be = getEmptyBackend()

	c3 := NC(&C{
		Data: map[string]interface{}{
			"type": 0,
		},
		Directives: Directives{
			HTML: HTMLDirective{
				Render: func(this *Component) string {
					return `<a href="https://gascore.github.io/" target="_blank">project site</a>`
				},
			},
		},
	}, func(this *Component) []interface{} {
		return []interface{}{}
	})
	c3.be = getEmptyBackend()

	c4 := NC(&C{
		Data: map[string]interface{}{
			"hello_world": "Hello world!!",
			"bye_world":   "no",
		},
		be: getEmptyBackend(),
	}, func(this *Component) []interface{} {
		return []interface{}{
			NE(&C{
				Binds: map[string]Bind{
					"color": func() string {
						if this.GetData("bye_world").(string) == "no" {
							return "hello"
						}
						return "bye"
					},
				},
			}, "lorem ipsum"),
		}
	})
	c4.be = getEmptyBackend()

	data := []struct {
		c     *Component
		value interface{}
		query string
		error error
		extra func(*C) error
	}{
		{
			c:     c1,
			query: "hello_world",
			value: "Hello world!",
			error: nil,
		},
		{
			c:     c1,
			query: "bye_world",
			value: "bye",
			error: errInByeWatcher,
		},
		{
			c: &C{
				Data: map[string]interface{}{},
			},
			query: "hello_world",
			value: "Hello world!",
			error: ErrNilField,
		},
		{
			c:     &C{}, // check will GetData throw panic with nil Data map
			query: "hello_world",
			value: "...",
			error: ErrComponentDataIsNil,
		},
		{ // html directive changed
			c:     c2,
			query: "type",
			value: 1,
			error: nil,
		},
		{ // html directive not changed
			c:     c3,
			query: "type",
			value: 1,
			error: nil,
		},
		{ // RenderTree
			c:     c4,
			query: "hello_world",
			value: "Hello world!",
			error: nil,
			extra: func(this *C) error {
				if len(this.RChildes) == 0 {
					return errors.New("component childes are nil")
				}

				if this.RChildes[0].(*Component).RenderedBinds["color"] != "hello" {
					return errors.New("invalid child rendered bind")
				}
				return nil
			},
		},
		{ // RenderTree
			c:     c4,
			query: "bye_world",
			value: "Goodbye!",
			error: nil,
			extra: func(this *C) error {
				if len(this.RChildes) == 0 {
					return errors.New("component childes are nil")
				}

				if this.RChildes[0].(*Component).RenderedBinds["color"] != "bye" {
					return errors.New("invalid child rendered bind")
				}
				return nil
			},
		},
	}

	for _, el := range data {
		err := el.c.SetData(el.query, el.value)
		if err != nil && (el.error == nil || err != el.error) {
			t.Errorf("unexpected error: %s", err.Error())
			continue
		}

		if el.error == nil && el.extra != nil {
			err := el.extra(el.c)
			if err != nil {
				t.Errorf("error in extra test: %s", err.Error())
				continue
			}
		}
	}
}

func TestRemove(t *testing.T) {
	// test remove function
	a := remove([]interface{}{1, 2, 3}, 1)
	aRes := []interface{}{1, 3}
	if len(a) != len(aRes) || a[0] != aRes[0] || a[1] != aRes[1] {
		t.Errorf("invalid remove function response: want: %v, got: %v", aRes, a)
	}
}

func TestDataDeleteFromArray(t *testing.T) {
	c1 := NC(
		&C{
			Data: map[string]interface{}{
				"arr":  []interface{}{1, 2, 3},
				"arr2": []int{1, 2, 3},
			},
		},
		func(this *Component) []interface{} {
			return []interface{}{len(this.GetData("arr").([]interface{}))}
		})
	c1.be = getEmptyBackend()

	data := []struct {
		c     *C
		query string
		index int
		len   int
		err   error
	}{
		{
			c:     c1,
			query: "arr",
			index: 1,
			len:   2,
			err:   nil,
		},
		{
			c:     c1,
			query: "arr2",
			index: 1,
			err:   ErrInvalidDataFieldType,
		},
	}

	for _, el := range data {
		err := el.c.DataDeleteFromArray(el.query, el.index)
		if err != nil && (el.err == nil || err != el.err) {
			t.Errorf("unexpected error: %s", err.Error())
			continue
		}

		if el.err == nil && len(el.c.GetData(el.query).([]interface{})) != el.len {
			t.Errorf("invalid array length, want: %d, got: %d", el.len, len(el.c.GetData(el.query).([]interface{})))
			continue
		}
	}
}

func TestDataAddToArray(t *testing.T) {
	c1 := NC(
		&C{
			Data: map[string]interface{}{
				"arr":  []interface{}{1, 2, 3},
				"arr2": []int{1, 2, 3},
			},
		},
		func(this *Component) []interface{} {
			return []interface{}{len(this.GetData("arr").([]interface{}))}
		})
	c1.be = getEmptyBackend()

	data := []struct {
		c     *C
		query string
		value interface{}
		len   int
		err   error
	}{
		{
			c:     c1,
			query: "arr",
			value: 4,
			len:   4,
			err:   nil,
		},
		{
			c:     c1,
			query: "arr2",
			value: 4,
			len:   4,
			err:   ErrInvalidDataFieldType,
		},
	}

	for _, el := range data {
		err := el.c.DataAddToArray(el.query, el.value)
		if err != nil && (el.err == nil || err != el.err) {
			t.Errorf("unexpected error: %s", err.Error())
			continue
		}

		if el.err == nil && len(el.c.GetData(el.query).([]interface{})) != el.len {
			t.Errorf("invalid array length, want: %d, got: %d", el.len, len(el.c.GetData(el.query).([]interface{})))
			continue
		}
	}
}

func TestDataEditArray(t *testing.T) {
	c1 := NC(
		&C{
			Data: map[string]interface{}{
				"arr":  []interface{}{1, 2, 3},
				"arr2": []int{1, 2, 3},
			},
		},
		func(this *Component) []interface{} {
			return []interface{}{len(this.GetData("arr").([]interface{}))}
		})
	c1.be = getEmptyBackend()

	data := []struct {
		c     *C
		query string
		index int
		value interface{}
		len   int
		err   error
	}{
		{
			c:     c1,
			query: "arr",
			index: 1,
			value: 22,
			len:   3,
			err:   nil,
		},
		{
			c:     c1,
			query: "arr2",
			index: 1,
			value: 22,
			len:   3,
			err:   ErrInvalidDataFieldType,
		},
	}

	for _, el := range data {
		err := el.c.DataEditArray(el.query, el.index, el.value)
		if err != nil && (el.err == nil || err != el.err) {
			t.Errorf("unexpected error: %s", err.Error())
			continue
		}

		if el.err == nil && len(el.c.GetData(el.query).([]interface{})) != el.len {
			t.Errorf("invalid array length, want: %d, got: %d", el.len, len(el.c.GetData(el.query).([]interface{})))
			continue
		}
	}
}

func TestDataDeleteFromMap(t *testing.T) {
	c1 := NC(
		&C{
			Data: map[string]interface{}{
				"arr": map[string]interface{}{
					"foo": 1,
					"bar": 2,
				},
				"arr2": map[string]int{
					"foo": 1,
					"bar": 2,
				},
			},
		},
		func(this *Component) []interface{} {
			return []interface{}{len(this.GetData("arr").(map[string]interface{}))}
		})
	c1.be = getEmptyBackend()

	data := []struct {
		c     *C
		query string
		key   string
		len   int
		err   error
	}{
		{
			c:     c1,
			query: "arr",
			key:   "bar",
			len:   1,
			err:   nil,
		},
		{
			c:     c1,
			query: "arr2",
			key:   "bar",
			len:   1,
			err:   ErrInvalidDataFieldType,
		},
	}

	for _, el := range data {
		err := el.c.DataDeleteFromMap(el.query, el.key)
		if err != nil && (el.err == nil || err != el.err) {
			t.Errorf("unexpected error: %s", err.Error())
			continue
		}

		if el.err == nil && len(el.c.GetData(el.query).(map[string]interface{})) != el.len {
			t.Errorf("invalid map length, want: %d, got: %d", el.len, len(el.c.GetData(el.query).(map[string]interface{})))
			continue
		}
	}
}

func TestDataEditMap(t *testing.T) {
	c1 := NC(
		&C{
			Data: map[string]interface{}{
				"arr": map[string]interface{}{
					"foo": 1,
					"bar": 2,
				},
				"arr2": map[string]int{
					"foo": 1,
					"bar": 2,
				},
			},
		},
		func(this *Component) []interface{} {
			return []interface{}{len(this.GetData("arr").(map[string]interface{}))}
		})
	c1.be = getEmptyBackend()

	data := []struct {
		c     *C
		query string
		key   string
		value interface{}
		len   int
		err   error
	}{
		{
			c:     c1,
			query: "arr",
			key:   "bar",
			value: 22,
			len:   2,
			err:   nil,
		},
		{
			c:     c1,
			query: "arr2",
			key:   "bar",
			value: 22,
			len:   2,
			err:   ErrInvalidDataFieldType,
		},
	}

	for _, el := range data {
		err := el.c.DataEditMap(el.query, el.key, el.value)
		if err != nil && (el.err == nil || err != el.err) {
			t.Errorf("unexpected error: %s", err.Error())
			continue
		}

		if el.err == nil && len(el.c.GetData(el.query).(map[string]interface{})) != el.len {
			t.Errorf("invalid map length, want: %d, got: %d", el.len, len(el.c.GetData(el.query).(map[string]interface{})))
			continue
		}
	}
}
