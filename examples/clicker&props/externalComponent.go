package main

import (
	"fmt"
	"github.com/Sinicablyat/gas/core"
)

// GetNumberViewer return very cool number viewer.
// It can be in another directory too.
// For reference from not parent component you can use `values` (they will reload).
func GetNumberViewer(this *core.Component, values ...interface{}) interface{} {
	return core.NewComponent(
		&core.Component{
			ParentC: this,
			Tag: "i",
			Attrs: map[string]string{
				"id": "needful_wrapper--number-viewer",
			},
		},
		func(this3 *core.Component) interface{} {
			return fmt.Sprintf("%d times", values[0])
		})
}
