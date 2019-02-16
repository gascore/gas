package gas

// Eb empty backend for testing backend calling methods
type Eb struct{}

func getEmptyBackend() BackEnd {
	return Eb{}
}

func (e Eb) New(a string) (string, error) {
	return "app", nil
}

func (e Eb) Init(g Gas) error {
	return nil
}

func (e Eb) UpdateComponentChildes(c *Component, newChildes, oldChildes []interface{}) error {
	// return errors.New("not supported")
	return nil
}

func (e Eb) ReCreate(c *Component) error {
	// return errors.New("not supported")
	return nil
}

func (e Eb) GetElement(c *Component) interface{} {
	if c.Attrs != nil && c.Attrs["need-component"] == "true" {
		return "not nil!"
	}

	return nil
}

func (e Eb) GetGasEl(g *Gas) interface{} {
	return "not nil!"
}

func (e Eb) ConsoleLog(values ...interface{}) {}

func (e Eb) ConsoleError(values ...interface{}) {}
