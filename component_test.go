package gas

import (
	"testing"
)

type TestNewForData struct {
	data          string
	renderer      func(int, interface{}) interface{}
	childesLength int
}

func TestNewComponent(t *testing.T) {
	data := []*Component{
		NewComponent(
			&C{
				Tag: "h1",
				Attrs: map[string]string{
					"id": "foo",
				},
			},
			func(this *Component) []interface{} {
				return ToGetComponentList(
					"Some ",
					NE(&C{Tag: "b"}, "bold"),
					"text")
			},
		),
		NewComponent(
			&C{
				Attrs: map[string]string{
					"id": "bar",
				},
			},
			func(this *Component) []interface{} {
				return ToGetComponentList(
					"Some ",
					NE(&C{Tag: "i"}, "italic"),
					"text")
			},
		),
		NewBasicComponent(
			&C{
				Attrs: map[string]string{
					"id": "bar",
				},
			},
			func(this *Component) []interface{} {
				return ToGetComponentList(
					"Some ",
					NE(&C{Tag: "i"}, "italic"),
					"text")
			},
		),
		NewComponent(
			&C{
				Attrs: map[string]string{
					"id": "hide",
				},
			},
			func(this *Component) []interface{} {
				return ToGetComponentList(
					NE(
						&C{
							If: func(p *Component) bool {
								return this.Attrs["id"] != "hide"
							},
						},
						"wow"),
					"another one")
			},
		),
	}

	for _, c := range data {
		if c.UUID == "" {
			t.Errorf("invalid comoponent UUID: lenght < 0")
		}

		if c.Tag == "" {
			t.Errorf("invalid comoponent tag: lenght < 0")
		}

		if c.isElement != c.IsElement() {
			t.Errorf("wtf??!?!?")
		}

		childes := c.Childes(c)
		if len(childes) == 0 {
			t.Errorf("components childes are null")
		}
	}
}

type TestGetElementData struct {
	c     *Component
	isNil bool
}

func TestGetElement(t *testing.T) {
	data := []TestGetElementData{
		{
			c: NC(
				&C{
					Tag: "div",
				},
				func(this *Component) []interface{} {
					return []interface{}{"nil"}
				},
			),
			isNil: true,
		},
		{
			c: NC(
				&C{
					Tag: "div",
					Attrs: map[string]string{
						"need-component": "true",
					},
				},
				func(this *Component) []interface{} {
					return []interface{}{"nil"}
				},
			),
			isNil: false,
		},
	}

	for _, el := range data {
		el.c.RC = GetEmptyRenderCore()
		_c := el.c.Element()
		if _c == nil && !el.isNil {
			t.Error("meh, I just want 100% coverage")
			continue
		}

		_c = el.c.GetElementUnsafely()
		if _c == nil && !el.isNil {
			t.Error("meh, I just want 100% coverage")
		}
	}
}

func TestNewFor(t *testing.T) {
	this := &C{
		Data: map[string]interface{}{
			"arr": []interface{}{
				"foo",
				"bar",
				"lol",
			},
		},
		Attrs: map[string]string{
			"id": "list",
		},
		RC: GetEmptyRenderCore(),
	}

	data := []TestNewForData{
		{
			data: "arr",
			renderer: func(i int, el interface{}) interface{} {
				return NE(
					&C{
						Tag: "li",
						Attrs: map[string]string{
							"class": "list__item",
						},
					},
					i, ": ", el,
				)
			},
			childesLength: 3,
		},
		{
			data: "anotherarr", // error here
			renderer: func(i int, el interface{}) interface{} {
				return NE(
					&C{
						Tag: "li",
						Attrs: map[string]string{
							"class": "list__item",
						},
					},
					i, ": ", el,
				)
			},
			childesLength: 0,
		},
		{
			data: "arr",
			renderer: func(i int, el interface{}) interface{} {
				return NE(
					&C{
						Tag: "li",
					},
					i, ": ", el,
				)
			},
			childesLength: 3,
		},
	}

	for _, el := range data {
		childesRenderer := func(this *Component) []interface{} {
			return ToGetComponentList(NE(&C{Tag: "ul"}, NewFor(el.data, this, el.renderer)))
		}

		c := NC(this, childesRenderer)
		c.RC = GetEmptyRenderCore()

		elementC := I2C(c.Childes(c)[0])
		elementC.RC = GetEmptyRenderCore()

		childes := elementC.Childes(elementC)

		if len(childes) != el.childesLength {
			t.Errorf("invalid childes list length got: %d, want: %d", len(childes), el.childesLength)
		}
	}
}

func TestIsString(t *testing.T) {
	s := "string"
	if !IsString(s) {
		t.Errorf("wtf?!?!: this string not string: %s", s)
	}
}
