package gas

// Hooks component lifecycle hooks
type Hooks struct {
	Created       Hook // When component has been created in golang only (GetElement isn't available)
	Mounted       Hook // When component has been mounted (GetElement is available)
	WillDestroy   Hook // Before component destroy (GetElement is available)
	BeforeUpdate  Hook // When component child don't updated
	Updated		  Hook // After component child was updated
}

// Hook - lifecycle hook
type Hook func(*Component) error


func RunMountedIfCan(i interface{}) error {
	if !IsComponent(i) {
		return nil
	}

	c := I2C(i)

	for _, child := range c.RChildes {
		if !IsComponent(child) {
			continue
		}

		err := RunMountedIfCan(I2C(child))
		if err != nil {
			return err
		}
	}

	if c.Hooks.Mounted != nil {
		err := c.Hooks.Mounted(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func RunWillDestroyIfCan(i interface{}) error {
	if !IsComponent(i) {
		return nil
	}

	c := I2C(i)
	for _, child := range c.RChildes {
		if !IsComponent(child) {
			continue
		}

		err := RunWillDestroyIfCan(I2C(child))
		if err != nil {
			return err
		}
	}

	if c.Hooks.WillDestroy != nil {
		err := c.Hooks.WillDestroy(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func RunUpdatedIfCan(i interface{}) error {
	if !IsComponent(i) {
		return nil
	}

	// run Updated hook for component parent(!)
	c := I2C(i).Parent

	if c.Hooks.Updated != nil {
		err := c.Hooks.Updated(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func RunBeforeUpdateIfCan(i interface{}) error {
	if !IsComponent(i) {
		return nil
	}

	// run BeforeUpdate hook for component parent(!)
	c := I2C(i).Parent

	if c.Hooks.BeforeUpdate != nil {
		err := c.Hooks.BeforeUpdate(c)
		if err != nil {
			return err
		}
	}

	return nil
}
