package gas

import (
	"fmt"
	"reflect"
)

// Changed return true if node Changed
func Changed(new, old interface{}) (bool, error) {
	if fmt.Sprintf("%T", new) != fmt.Sprintf("%T", old) {
		return true, nil
	}

	if IsString(new) {
		return new.(string) != old.(string), nil
	} else if IsComponent(new) {
		return !isComponentsEquals(I2C(new), I2C(old)), nil
	}

	return false, fmt.Errorf("changed: invalid `new` or `old`. types: %T, %T", new, old)
}

func isComponentsEquals(new, old *Component) bool {
	return  new.Tag == old.Tag &&
			new.Directives.HTML.Rendered == old.Directives.HTML.Rendered &&

			reflect.DeepEqual(new.Attrs, old.Attrs) &&
			reflect.DeepEqual(new.RenderedBinds, old.RenderedBinds) &&

			compareHooks(new, old) &&

			reflect.ValueOf(new.Directives.HTML.Render).Pointer() == reflect.ValueOf(old.Directives.HTML.Render).Pointer() &&
			reflect.ValueOf(new.Directives.If).Pointer() == reflect.ValueOf(old.Directives.If).Pointer() &&
			compareForDirectives(new, old) &&
			new.Directives.Model.Data == old.Directives.Model.Data
}

func compareHooks(new, old *Component) bool {
	if !new.isElement && !old.isElement {
		return true
	}
	
	if (new.isElement && !old.isElement) || (!new.isElement && old.isElement) {
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