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

// CallBeforeCreatedIfCan call component and it's childes BeforeCreated hook
func CallBeforeCreatedIfCan(i interface{}) error {
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
			e.RChildes = e.RenderTree()
		}
	}

	for _, child := range e.RChildes {
		err := CallBeforeCreatedIfCan(child)
		if err != nil {
			return err
		}
	}

	return nil
}

// CallMountedIfCan call component and it's childes Mounted hook
func CallMountedIfCan(i interface{}) error {
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

	for _, child := range e.RChildes {
		err := CallMountedIfCan(child)
		if err != nil {
			return err
		}
	}

	return nil
}

// CallBeforeDestroyIfCan call component and it's childes WillDestroy hook
func CallBeforeDestroyIfCan(i interface{}) error {
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

	for _, child := range e.RChildes {
		err := CallBeforeDestroyIfCan(child)
		if err != nil {
			return err
		}
	}

	return nil
}

// CallUpdatedIfCan call component parent (true component) Updated hook
func CallUpdatedIfCan(i interface{}) error {
	e, ok := i.(*Element)
	if !ok {
		return nil
	}

	// run Updated hook for component parent
	c := e.ParentComponent().Component

	if c.Hooks.Updated != nil {
		err := c.Hooks.Updated()
		if err != nil {
			return err
		}
	}

	return nil
}

// CallBeforeUpdateIfCan call component parent (true component) BeforeUpdate
func CallBeforeUpdateIfCan(i interface{}) error {
	e, ok := i.(*Element)
	if !ok {
		return nil
	}

	// run Updated hook for component parent
	c := e.ParentComponent().Component

	if c.Hooks.BeforeUpdate != nil {
		err := c.Hooks.BeforeUpdate()
		if err != nil {
			return err
		}
	}

	return nil
}
