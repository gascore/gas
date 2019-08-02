package gas

// Hooks component lifecycle hooks
type Hooks struct {
	BeforeCreated HookWithControl // When parent already rendered (appended to DOM), but component Element don't yet (you can rerender childes)
	Created       Hook            // When component has been created in golang only (Element isn't available)

	Mounted Hook // When component has been mounted (Element is available)

	BeforeDestroy Hook // Before component destroy (Element is available)

	BeforeUpdate Hook // When component child don't updated
	Updated      Hook // After component child was updated
}

// Hook - lifecycle hook
type Hook func() error

// HookWithControl - lifecycle hook. Return true for rerender component childes
type HookWithControl func() (rerender bool, err error)

// CallBeforeCreated call component and it's childes BeforeCreated hook
func CallBeforeCreated(i interface{}) error {
	e, ok := i.(*Element)
	if !ok {
		return nil
	}

	c := e.Component
	if c != nil && c.Hooks.BeforeCreated != nil {
		rerender, err := c.Hooks.BeforeCreated()
		if err != nil {
			return err
		}

		if rerender {
			e.UpdateChildes()
		}
	}

	for _, child := range e.Childes {
		err := CallBeforeCreated(child)
		if err != nil {
			return err
		}
	}

	return nil
}

// CallMounted call component and it's childes Mounted hook
func CallMounted(i interface{}) error {
	e, ok := i.(*Element)
	if !ok {
		return nil
	}

	c := e.Component
	if c != nil && c.Hooks.Mounted != nil {
		err := c.Hooks.Mounted()
		if err != nil {
			return err
		}
	}

	for _, child := range e.Childes {
		err := CallMounted(child)
		if err != nil {
			return err
		}
	}

	return nil
}

// CallBeforeDestroy call component and it's childes WillDestroy hook
func CallBeforeDestroy(i interface{}) error {
	e, ok := i.(*Element)
	if !ok {
		return nil
	}

	c := e.Component
	if c != nil && c.Hooks.BeforeDestroy != nil {
		err := c.Hooks.BeforeDestroy()
		if err != nil {
			return err
		}
	}

	for _, child := range e.Childes {
		err := CallBeforeDestroy(child)
		if err != nil {
			return err
		}
	}

	return nil
}

// CallUpdated call Updated hook
func CallUpdated(p *Element) error {
	c := p.Component
	if c == nil {
		c = p.ParentComponent().Component
	}

	if c.Hooks.Updated != nil {
		err := c.Hooks.Updated()
		if err != nil {
			return err
		}
	}

	return nil
}

// CallBeforeUpdate call BeforeUpdate hook
func CallBeforeUpdate(element *Element) error {
	c := element.Component
	if c == nil {
		c = element.ParentComponent().Component
	}

	if c.Hooks.BeforeUpdate != nil {
		err := c.Hooks.BeforeUpdate()
		if err != nil {
			return err
		}
	}

	return nil
}
