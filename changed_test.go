package gas

import (
	"fmt"
	"log"
	"testing"
)

type TestChangedData struct {
	x         interface{}
	y         interface{}
	answer    bool
	haveError bool
}

func TestChanged(t *testing.T) {
	c1 := &C{
		Data: map[string]interface{}{
			"1": 2,
		},
		Hooks: Hooks{
			Mounted: func(this *Component) error {
				fmt.Println("it is interesting")
				return nil
			},
		},
		Tag: "h1",
	}
	c2 := &C{
		Data: map[string]interface{}{
			"1": 2,
		},
		Hooks: Hooks{
			Created: func(this *Component) error {
				fmt.Println("it is not interesting")
				return nil
			},
		},
		Methods: map[string]Method{
			"1": func(this *Component, values ...interface{}) error {
				log.Println("method 1")
				return nil
			},
		},
		Tag: "h2",
	}
	c3 := &C{
		Data: map[string]interface{}{
			"1":    2,
			"such": "empty",
		},
		Hooks: Hooks{
			Mounted: func(this *Component) error {
				fmt.Println("it is interesting")
				return nil
			},
		},
		Tag: "h1",
	}
	c4 := &C{
		Hooks: Hooks{
			Created: func(this *Component) error {
				fmt.Println("it is not interesting")
				return nil
			},
		},
		Methods: map[string]Method{
			"1": func(this *Component, values ...interface{}) error {
				log.Println("method 1")
				return nil
			},
			"2": func(this *Component, values ...interface{}) error {
				log.Println("method 2")
				return nil
			},
		},
		Tag: "h2",
	}
	c5 := &C{
		Hooks: Hooks{
			Created: func(this *Component) error {
				fmt.Println("it is not interesting")
				return nil
			},
		},
		Methods: map[string]Method{
			"1": func(this *Component, values ...interface{}) error {
				log.Println("method 1")
				return nil
			},
			"2": func(this *Component, values ...interface{}) error {
				log.Println("method 31231231231231231313131")
				return nil
			},
		},
		Tag: "h2",
	}
	c6 := &C{
		Hooks: Hooks{
			Created: func(this *Component) error {
				fmt.Println("it is not interesting")
				return nil
			},
		},
		Watchers: map[string]Watcher{
			"1": func(this *Component, new, old interface{}) error {
				log.Println("watcher 1")
				return nil
			},
			"2": func(this *Component, new, old interface{}) error {
				log.Println("watcher2 2")
				return nil
			},
		},
		Tag: "h2",
	}
	c7 := &C{
		Hooks: Hooks{
			Created: func(this *Component) error {
				fmt.Println("it is not interesting")
				return nil
			},
		},
		Watchers: map[string]Watcher{
			"1": func(this *Component, new, old interface{}) error {
				log.Println("watcher 1")
				return nil
			},
			"2": func(this *Component, new, old interface{}) error {
				log.Println("watcher 31231231231231231313131")
				return nil
			},
		},
		Tag: "h2",
	}
	c8 := &C{
		Hooks: Hooks{
			Created: func(this *Component) error {
				fmt.Println("it is not interesting")
				return nil
			},
		},
		Computeds: map[string]Computed{
			"1": func(this *Component, values ...interface{}) (interface{}, error) {
				log.Println("computed 1")
				return 1, nil
			},
			"2": func(this *Component, values ...interface{}) (interface{}, error) {
				log.Println("computed 2")
				return 2, nil
			},
		},
		Tag: "h2",
	}
	c9 := &C{
		Hooks: Hooks{
			Created: func(this *Component) error {
				fmt.Println("it is not interesting")
				return nil
			},
		},
		Computeds: map[string]Computed{
			"1": func(this *Component, values ...interface{}) (interface{}, error) {
				log.Println("watcher 1")
				return 1, nil
			},
			"2": func(this *Component, values ...interface{}) (interface{}, error) {
				log.Println("watcher 31231231231231231313131")
				return 33243242, nil
			},
		},
		Tag: "h2",
	}
	c10 := &C{
		Hooks: Hooks{
			Created: func(this *Component) error {
				fmt.Println("it is not interesting")
				return nil
			},
		},
		Computeds: map[string]Computed{
			"1": func(this *Component, values ...interface{}) (interface{}, error) {
				log.Println("watcher 1")
				return 1, nil
			},
			"2": func(this *Component, values ...interface{}) (interface{}, error) {
				log.Println("watcher 231231231231231313131")
				return 243242, nil
			},
			"3": func(this *Component, values ...interface{}) (interface{}, error) {
				log.Println("watcher 3")
				return 3, nil
			},
		},
		Tag: "h2",
	}
	c11 := &C{
		Hooks: Hooks{
			Created: func(this *Component) error {
				fmt.Println("it is not interesting")
				return nil
			},
		},
		Tag: "h2",
	}

	e1 := &C{
		Tag:       "h1",
		isElement: true,
	}
	e2 := &C{
		Attrs: map[string]string{
			"class": "mega class",
		},
		Tag:       "h2",
		isElement: true,
	}
	e3 := &C{
		Tag:       "h1",
		isElement: false,
	}
	e4 := &C{
		Tag:       "li",
		isElement: true,
		Directives: Directives{
			For: ForDirective{
				isItem:       true,
				itemValueI:   0,
				itemValueVal: "some",
			},
		},
	}
	e5 := &C{
		Tag:       "li",
		isElement: true,
	}
	e6 := &C{
		Tag:       "li",
		isElement: true,
		Directives: Directives{
			For: ForDirective{
				isItem:       true,
				itemValueI:   1,
				itemValueVal: "some",
			},
		},
	}

	data := []TestChangedData{
		{0, 1, true, false},
		{1, 1, false, false},
		{"test", "test1", true, false},
		{"test", "test", false, false},

		{c1, c2, true, false},
		{c1, c1, false, false},
		{c3, c2, true, false},
		{c1, c3, true, false},             // c1 and c3 not same because CompareHooks check hooks pointers, not what they are doing
		{c2, c4, true, false},             // methods
		{c4, c5, true, false},             // methods
		{c6, c7, true, false},             // watchers
		{c8, c9, true, false},             // computeds
		{c9, c10, true, false},            // computeds
		{c11, &C{Tag: "h2"}, true, false}, // hooks

		{e1, e2, true, false},
		{e1, e1, false, false},
		{e3, e2, true, false},
		{e1, e3, true, false},
		{e4, e5, true, false},
		{e4, e6, true, false},

		{0, "1", true, false},
		{[]int{1, 2}, []int{3, 4}, false, true},
	}

	for i, el := range data {
		isChanged, err := Changed(el.x, el.y)
		if err != nil {
			if el.haveError {
				continue
			}

			t.Errorf("error while testing Changed: %s", err.Error())
			return
		}

		if isChanged != el.answer {
			t.Errorf("Compare of %v and %v was incorrect, got: %t, want: %t (i: %d)", el.x, el.y, isChanged, el.answer, i)
		}
	}
}
