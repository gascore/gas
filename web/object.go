package web

import (
	"fmt"
	"strconv"

	"github.com/gascore/dom"
	"github.com/gascore/gas"
	"syscall/js"
)

func parseInt(a string) (int, error) {
	if a == "" {
		return 0, nil
	}

	return strconv.Atoi(a)
}

func warnError(err error) {
	dom.ConsoleError(err.Error())
}

type object struct{ o js.Value }

func (o object) String() string                   { return o.o.String() }
func (o object) Int() int                         { return o.o.Int() }
func (o object) Float() float64                   { return o.o.Float() }
func (o object) Get(q string) gas.Object          { return object{o: o.o.Get(q)} }
func (o object) Set(name string, val interface{}) { o.o.Set(name, val) }
func (o object) GetString(q string) string        { return o.o.Get(q).String() }
func (o object) GetBool(q string) bool            { return o.o.Get(q).Bool() }
func (o object) GetInt(q string) int              { return o.o.Get(q).Int() }
func (o object) Raw() interface{}                 { return o.o }
func (o object) Call(name string, args ...interface{}) gas.Object {
	return object{o: o.o.Call(name, args...)}
}

// ToUniteObject convert dom.Value to gas.Object
func ToUniteObject(e dom.Value) gas.Object { return object{o: e.JSValue()} }

type event struct{ 
	gas.Object
	event dom.Event
	isCheckbox bool
}

// Value reurn event value
func (e event) Value() string {
	if e.isCheckbox {
		if e.ValueBool() {
			return "true"
		}
		return "false"
	}

	return e.event.Target().Value()
}

// Value reurn event value and convert it to int
func (e event) ValueInt() int {
	val := e.Value()
	n, err := strconv.Atoi(val)
	if err != nil {
		dom.ConsoleError(fmt.Sprintf("cannot convert event value to int: \"%s\"", val))
	}

	return n
}

// Value reurn event value and convert it to boolean
func (e event) ValueBool() bool {
	if e.isCheckbox {
		return e.event.Target().JSValue().Get("checked").Bool()
	}

	val := e.Value()
	if val == "true" {
		return true
	} else if val != "false" {
		dom.ConsoleError(fmt.Sprintf("cannot convert event value to bool: \"%s\"", val))
	}

	return false
}

// ToGasEvent convert dom.Event to gas.Event
func ToGasEvent(domEvent dom.Event, isCheckbox bool) gas.Event { 
	e := event{ToUniteObject(domEvent), domEvent, isCheckbox}
	return e
}
