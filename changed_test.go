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
	notGopherjs bool
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
			"1": func(this *Component, values ...interface{}) (interface{}, error) {
				log.Println("method 1")
				return nil, nil
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
			"1": func(this *Component, values ...interface{}) (interface{}, error) {
				log.Println("method 1")
				return nil, nil
			},
			"2": func(this *Component, values ...interface{}) (interface{}, error) {
				log.Println("method 2")
				return nil, nil
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
			"1": func(this *Component, values ...interface{}) (interface{}, error) {
				log.Println("method 1")
				return nil, nil
			},
			"2": func(this *Component, values ...interface{}) (interface{}, error) {
				log.Println("method 31231231231231231313131")
				return nil, nil
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
		{x: 0, y: 1, answer: true},
		{x: 1, y: 1, answer: false},
		{x: "test", y: "test1", answer: true},
		{x: "test", y: "test", answer: false},

		{x: c1, y: c2, answer: true},
		{x: c1, y: c1, answer: false},
		{x: c3, y: c2, answer: true},
		{x: c1, y: c3, answer: true}, // c1 and c3 not same because CompareHooks check hooks pointers, not what they are doing. But not in gopherjs
		{x: c2, y: c4, answer: true, notGopherjs: true}, // methods
		{x: c4, y: c5, answer: true, notGopherjs: true}, // methods
		{x: c6, y: c7, answer: true, notGopherjs: true}, // watchers
		{x: c8, y: c9, answer: true, notGopherjs: true}, // computeds
		{x: c9, y: c10, answer: true}, // computeds
		{x: c11, y: &C{Tag: "h2"}, answer: true}, // hooks

		{x: e1, y: e2, answer:true},
		{x: e1, y: e1, answer: false},
		{x: e3, y: e2, answer: true},
		{x: e1, y: e3, answer: true},
		{x: e4, y: e5, answer: true},
		{x: e4, y: e6, answer: true},

		{x: 0, y: "1",  answer:true},
		{x: []int{1, 2}, y: []int{3, 4}, answer: false, haveError: true},
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

		if el.notGopherjs {
			continue
		}

		if isChanged != el.answer {
			t.Errorf("Compare %d was incorrect, got: %t, want: %t", i, isChanged, el.answer)
		}
	}
}
