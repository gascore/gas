package gas

// FunctionalComponent wrapper for Component with react hooks (in gas maner)
type FunctionalComponent struct {
	C *C

	statesCounter  int
	states         []interface{}
	effectsCounter int
	effects        []Hook

	renderer func() []interface{}
}

// Init create *C from *F
func (f *FunctionalComponent) Init(notPointer bool, renderer func() []interface{}) *E {
	f.renderer = renderer

	c := &C{
		NotPointer: notPointer,
		Root:       f,
		Hooks: Hooks{
			Updated: f.runEffects,
			Mounted: f.runEffects,
		},
	}
	f.C = c

	return c.Init()
}

// UseState create new state value
func (root *FunctionalComponent) UseState(defaultVal interface{}) (func() interface{}, func(interface{})) {
	i := root.statesCounter
	root.statesCounter++

	if len(root.states)-1 < i {
		root.states = append(root.states, defaultVal)
	}

	getVal := func() interface{} {
		return root.states[i]
	}

	setVal := func(newVal interface{}) {
		root.states[i] = newVal
		go root.C.Update()
	}

	return getVal, setVal
}

// UseEffect add effect
func (root *FunctionalComponent) UseEffect(f Hook) {
	i := root.effectsCounter
	root.effectsCounter++

	if len(root.effects)-1 < i {
		root.effects = append(root.effects, f)
	}
}

// Render return functionalComponent childes
func (root *FunctionalComponent) Render() []interface{} {
	root.statesCounter = 0
	root.effectsCounter = 0
	return root.renderer()
}

func (root *FunctionalComponent) runEffects() error {
	for _, effect := range root.effects {
		err := effect()
		if err != nil {
			return err
		}
	}

	return nil
}
