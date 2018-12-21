package gas

import (
	"fmt"
	"github.com/Sinicablyat/gas/core"
	"github.com/Sinicablyat/gas/wasm"
	"github.com/frankenbeanies/uuid4"

	"github.com/Sinicablyat/dom"
)

var be core.BackEnd

// New create new gas application with custom backend
func New(backEnd core.BackEnd, startPoint string, components ...core.GetComponent) (core.Gas, error) {
	core.SetBackEnd(backEnd)
	be = backEnd

	tagName, err := be.New(startPoint)
	if err != nil {
		return core.Gas{}, err
	}

	mainComponent := core.NewComponent(
		&core.Component{
			Tag: tagName,
			Attrs: map[string]string{
				"id": startPoint,
				"data-main": "true",
			},
			UUID: uuid4.New().String(),
		}, components...)

	gas := core.Gas{App: *mainComponent, StartPoint: startPoint}

	return gas, nil
}

// NewWasm create new gas application with wasm backend
func NewWasm(startPoint string, components ...core.GetComponent) (core.Gas, error) {
	return New(wasm.GetBackEnd(), startPoint, components...)
}

// Init initialize gas application
func Init(gas core.Gas) error {
	err := be.Init(gas)
	if err != nil {
		return err
	}

	return nil
}

// ToGetComponentList return array by many parameters, because it's pretty
func ToGetComponentList(childes ...core.GetComponent) []core.GetComponent {
	return childes
}

// WarnError log error
func WarnError(err error) {
	if err == nil {
		return
	}

	dom.ConsoleError(err.Error())
}

// WarnIfNot console error if !ok
func WarnIfNot(ok bool) {
	if ok {
		return
	}

	dom.ConsoleError(fmt.Sprintf("invalid data type"))
}

var signal = make(chan int)

// KeepAlive keep alive runtime, without it application will stop (user won't be able to init events)
func KeepAlive() {
	for {
		<- signal
	}
}