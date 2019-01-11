package gas

// Hooks component lifecycle hooks
type Hooks struct {
	Created     Hook // When component has been created in golang only (GetElement isn't available)
	Mounted     Hook // When component has been mounted (GetElement is available)
	WillDestroy Hook // Before component destroy (GetElement is available)
	//BeforeUpdate  Hook // Will add in the future
	//Updated		Hook // Will add in the future
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
