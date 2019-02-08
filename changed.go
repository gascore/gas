package gas

import (
	"fmt"
	"github.com/Sinicablyat/gas"
	"reflect"
)

// Changed return true if node Changed
func Changed(new, old interface{}) (bool, error) {
	if fmt.Sprintf("%T", new) != fmt.Sprintf("%T", old) {
		return true, nil
	}

	switch new.(type) {
	case *gas.Component:
		return !isComponentsEquals(I2C(new), I2C(old)), nil

	case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return new != old, nil

	default:
		return false, fmt.Errorf("changed: invalid `new` or `old`. types: %T, %T", new, old)
	}
}

func isComponentsEquals(new, old *Component) bool {
	return  new.Tag == old.Tag &&
			new.Directives.HTML.Rendered == old.Directives.HTML.Rendered &&

			reflect.DeepEqual(new.Attrs, old.Attrs) &&
			reflect.DeepEqual(new.RenderedBinds, old.RenderedBinds) &&

			reflect.DeepEqual(new.Data, old.Data) &&
			compareMethods(new.Methods, old.Methods) &&
			compareComputeds(new.Computeds, old.Computeds) &&

			compareHooks(new, old) &&

			reflect.ValueOf(new.Directives.HTML.Render).Pointer() == reflect.ValueOf(old.Directives.HTML.Render).Pointer() &&
			reflect.ValueOf(new.Directives.If).Pointer() == reflect.ValueOf(old.Directives.If).Pointer() &&
			compareForDirectives(new, old) &&
			new.Directives.Model.Data == old.Directives.Model.Data
}

func compareMethods(new map[string]Method, old map[string]Method) bool {
	if len(new) != len(old) {
		return false
	}

	for i, el := range new {
		if reflect.ValueOf(el).Pointer() != reflect.ValueOf(old[i]).Pointer() {
			return false
		}
	}

	return true
}

func compareComputeds(new map[string]Computed, old map[string]Computed) bool {
	if len(new) != len(old) {
		return false
	}

	for i, el := range new {
		if reflect.ValueOf(el).Pointer() != reflect.ValueOf(old[i]).Pointer() {
			return false
		}
	}

	return true
}

func compareHooks(new, old *Component) bool {
	if  (new.isElement && !old.isElement) || (!new.isElement && old.isElement) {
		return false
	}

	return compareHook(new.Hooks.Created, old.Hooks.Created) &&
		compareHook(new.Hooks.Mounted, old.Hooks.Mounted) &&
		compareHook(new.Hooks.WillDestroy, old.Hooks.WillDestroy) &&
		compareHook(new.Hooks.BeforeUpdate, old.Hooks.BeforeUpdate) &&
		compareHook(new.Hooks.Updated, old.Hooks.Updated)
}

func compareHook(new, old Hook) bool {
	if new == nil && old == nil {
		return true
	}

	if (new == nil && old != nil) || (new != nil && old == nil) {
		return false
	}

	return reflect.ValueOf(new).Pointer() == reflect.ValueOf(old).Pointer()
}

func compareForDirectives(new, old *Component) bool {
	/*
		It's really bad way to fix bug with not-updated i, el in components Methods.
		We can only ForceUpdate methods, binds, e.t.c. and don't ForceUpdate 'body', but it will be in the future...
	*/

	newIsItem, newI, newVal := new.ForItemInfo()
	oldIsItem, oldI, oldVal := old.ForItemInfo()

	if newIsItem != oldIsItem {
		return false
	} else {
		return newI == oldI && reflect.DeepEqual(newVal, oldVal)
	}
}