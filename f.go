package gas

type FunctionalComponent struct {
	c *C

	notPointer bool
	statesCounter    int
	states     []interface{}
	
	effects []func()

	childes func()[]interface{}
}

func (root *FunctionalComponent) UseState(defaultVal interface{}) (func()interface{}, func(interface{})) {
	i := root.statesCounter
	root.statesCounter++

	if len(root.states)-1 < i {
		root.states = append(root.states, defaultVal)
	}

	getVal := func()interface{}{
		return root.states[i]
	}
	setVal := func(newVal interface{}) {
		root.states[i] = newVal
		go root.c.Update()
	}

	return getVal, setVal
}

func (root *FunctionalComponent) UseEffect(f func()) {
	root.effects = append(root.effects, f)
}

func (root *FunctionalComponent) Init(childes func()[]interface{}) *E {
	root.childes = childes

	runEffects := func() error {
		for _, effect := range root.effects {
			err := effect()
			if err != nil {
				return err
			}
		}
		return nil
	}

	c := &C{
		NotPointer: root.notPointer,
		Root:       root,
		Hooks: gas.Hooks{
			Created: runEffects,
			Mounted: runEffects,
			Updated: runEffects,
		},
	}
	root.c = c

	return c.Init()
}

func (root *FunctionalComponent) Render() []interface{} {
	root.statesCounter = 0
	return root.childes()
}

func NewFunctionalComponent(notPointer bool) *FunctionalComponent {
	/*
		Example:
		func FunctionalExample() *gas.E {
			f := gas.NFC(true)

			getCounter, setCounter := f.UseState(0)

			return f.Init(func()[]interface{} {return gas.CL(
				gas.NE(&gas.E{Tag: "button", Handlers: map[string]gas.Handler{"click": func(e gas.Object) { setCounter(getCounter().(int) + 1) }}}, "+"),
				getCounter(),
				gas.NE(&gas.E{Tag: "button", Handlers: map[string]gas.Handler{"click": func(e gas.Object) { setCounter(getCounter().(int) - 1) }}}, "-"),
			)})
		}
	*/

	return &FunctionalComponent{
		counter:    0,
		states:     []interface{}{},
		notPointer: notPointer,
	}
}
