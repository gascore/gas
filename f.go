package gas

// FunctionalComponent wrapper for Component with react hooks (in gas maner)
type FunctionalComponent struct {
	C *Component

	statesCounter   int
	states          []interface{}
	effectsCounter  int
	effects         []Effect
	cleanersCounter int
	cleaners        []func()

	renderer func() []interface{}
}

// Effect functional components effect
type Effect func() (func(), error)

// Init create *C from *F
func (f *FunctionalComponent) Init(notPointer bool, renderer func() []interface{}) *E {
	f.renderer = renderer

	c := &C{
		NotPointer: notPointer,
		Root:       f,
		Hooks: Hooks{
			Updated:       f.runEffects,
			Mounted:       f.runEffects,
			BeforeDestroy: f.runCleaners,
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
		root.C.Update()
	}

	return getVal, setVal
}

// UseEffect add effect
func (root *FunctionalComponent) UseEffect(effect Effect) {
	i := root.effectsCounter
	root.effectsCounter++

	if len(root.effects)-1 < i {
		root.effects = append(root.effects, effect)
	}
}

func (root *FunctionalComponent) addCleaner(cleaner func()) {
	i := root.cleanersCounter
	root.cleanersCounter++

	if len(root.cleaners)-1 < i {
		root.cleaners = append(root.cleaners, cleaner)
	}
}

func (root *FunctionalComponent) runEffects() error {
	for _, effect := range root.effects {
		cleaner, err := effect()
		if err != nil {
			return err
		}

		if cleaner != nil {
			root.addCleaner(cleaner)
		}
	}

	return nil
}

func (root *FunctionalComponent) runCleaners() error {
	println("CLEAR THIS")
	for _, cleaner := range root.cleaners {
		cleaner()
	}

	return nil
}

// Render return functionalComponent childes
func (root *FunctionalComponent) Render() []interface{} {
	root.statesCounter = 0
	root.effectsCounter = 0
	root.cleanersCounter = 0
	return root.renderer()
}
