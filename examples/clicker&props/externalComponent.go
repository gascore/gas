package main

import (
	"fmt"
	"github.com/gascore/gas"
)

// GetNumberViewer return very cool number viewer.
// It can be in another directory too.
// For reference from not parent component you can use `values` (they will reload).
func GetNumberViewer(click int) interface{} {
	return gas.NC(
		&gas.Component{
			Tag: "i",
			Attrs: map[string]string{
				"id": "needful_wrapper--number-viewer",
			},
		},
		func(this *gas.Component) []interface{} {
			return gas.ToGetComponentList(
				fmt.Sprintf("%d times", click))
		},)
}
