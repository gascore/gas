//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "T=bool,string,int,int8,int16,int32,int64,uint,uint8,uint16,uint32,uint64,float32,float64"
package gas

import "github.com/cheekybits/genny/generic"

type T generic.Type

func (root *FunctionalComponent) UseStateT(defaultVal T) (func()T, func(T)) {
	getVal, setVal := root.UseState(defaultVal)
	return func()T { return getVal().(T) }, func(newVal T) { setVal(newVal) }
}