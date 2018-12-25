package gas

import (
	"fmt"
	"github.com/frankenbeanies/uuid4"
)

var be BackEnd

// New create new gas application with custom backend
func New(backEnd BackEnd, startPoint string, components ...GetComponent) (Gas, error) {
	be = backEnd

	tagName, err := be.New(startPoint)
	if err != nil {
		return Gas{}, err
	}

	mainComponent := NewComponent(
		&Component{
			Tag: tagName,
			Attrs: map[string]string{
				"id": startPoint,
				"data-main": "true",
			},
			UUID: uuid4.New().String(),
		}, components...)

	gas := Gas{App: *mainComponent, StartPoint: startPoint}

	return gas, nil
}

// Init initialize gas application
func Init(gas Gas) error {
	err := be.Init(gas)
	if err != nil {
		return err
	}

	return nil
}

// ToGetComponentList return array by many parameters, because it's pretty
func ToGetComponentList(childes ...GetComponent) []GetComponent {
	return childes
}

// WarnError log error
func WarnError(err error) {
	if err == nil {
		return
	}

	be.ConsoleError(err.Error())
}

// WarnIfNot console error if !ok
func WarnIfNot(ok bool) {
	if ok {
		return
	}

	be.ConsoleError(fmt.Sprintf("invalid data type"))
}

func ConsoleLog(a ...interface{})   { be.ConsoleLog(a...) }
func ConsoleError(a ...interface{}) { be.ConsoleError(a...) }

var signal = make(chan int)

// KeepAlive keep alive runtime, without it application will stop (user won't be able to init events)
func KeepAlive() {
	for {
		<- signal
	}
}