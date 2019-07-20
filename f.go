package gas

type FunctionalComponent struct {
	c *C

	counter int
	states  map[int]interface{}

	childes []interface{}
}

func (root *FunctionalComponent) UseState(defaultVal interface{}) (interface{}, func(interface{})) {
	i := root.counter
	root.counter++

	setVal := func(newVal interface{}) {
		root.states[i] = newVal
		go root.c.Update()
	}

	if len(root.states)-1 < i {
		root.states = append(root.states, defaultVal)
	}

	return root.states[i], setVal
}

func (root *FunctionalComponent) Init(childes ...interface{}) *E {
	root.childes = childes
	c := &C{
		Root: root,
	}
	root.c = c

	return c.Init()
}

func (root *FunctionalComponent) Render() []interface{} {
	root.counter = 0
	return root.childes
}

func NewFunctionalComponent() *FunctionalComponent {
	/*
		Example:
		function SomeComponent() *gas.E {
			f := gas.NFC()

			msg, setMsg := f.UseState("empty message")
			msg = msg.(string)

			return f.Init(
				gas.NE(
					&E{Tag: "h1"},
					msg,
				),
			)
		}
	*/

	return &FunctionalComponent{
		counter: 0,
		states:  make(map[int]interface{}),
	}
}
