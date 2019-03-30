package gas

import (
	"fmt"
	"reflect"
)

// Changed return true if node Changed
func Changed(newEl, oldEl interface{}) (bool, error) {
	if fmt.Sprintf("%T", newEl) != fmt.Sprintf("%T", oldEl) {
		return true, nil
	}

	switch newEl.(type) {
	case *Component:
		return !isComponentsEquals(I2C(newEl), I2C(oldEl)), nil

	case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return newEl != oldEl, nil

	default:
		return false, fmt.Errorf("changed: invalid `newEl` or `oldEl`. types: %T, %T", newEl, oldEl)
	}
}

func isComponentsEquals(newC, oldC *Component) bool {
	if newC.isElement != oldC.isElement {
		return false
	}

	// if components are elements
	if newC.isElement && oldC.isElement {
		return newC.Tag == oldC.Tag &&
			newC.Ref == oldC.Ref &&
			newC.HTML.Rendered == oldC.HTML.Rendered &&

			compareAttributes(newC.Attrs, oldC.Attrs) &&
			compareAttributes(newC.RenderedBinds, oldC.RenderedBinds) &&

			reflect.ValueOf(newC.HTML.Render).Pointer() == reflect.ValueOf(oldC.HTML.Render).Pointer() &&
			reflect.ValueOf(newC.If).Pointer() == reflect.ValueOf(oldC.If).Pointer() &&
			compareForDirectives(newC, oldC) &&
			newC.Model.Data == oldC.Model.Data
	}

	// if components are *true* components
	return newC.Tag == oldC.Tag &&
		newC.HTML.Rendered == oldC.HTML.Rendered &&
		newC.RefsAllowed == oldC.RefsAllowed &&

		// reflect.DeepEqual(newC.Data, oldC.Data) &&
		compareWatchers(newC.Watchers, oldC.Watchers) &&
		compareMethods(newC.Methods, oldC.Methods) &&

		compareHooks(newC, oldC)
}

func compareAttributes(newMap, oldMap map[string]string) bool {
	if len(newMap) != len(oldMap) {
		return false
	}

	for i, el := range newMap {
		if el != oldMap[i] {
			return false
		}
	}

	return true
}

func compareWatchers(a, b map[string]Watcher) bool {
	n := make(map[string]interface{})
	for key, val := range a {
		n[key] = val
	}

	o := make(map[string]interface{})
	for key, val := range b {
		o[key] = val
	}

	return compareMapStringFunc(n, o)
}

func compareMethods(a, b map[string]Method) bool {
	n := make(map[string]interface{})
	for key, val := range a {
		n[key] = val
	}

	o := make(map[string]interface{})
	for key, val := range b {
		o[key] = val
	}

	return compareMapStringFunc(n, o)
}

func compareMapStringFunc(newMap, oldMap map[string]interface{}) bool {
	if len(newMap) != len(oldMap) {
		return false
	}

	for i, el := range newMap {
		if reflect.ValueOf(el).Pointer() != reflect.ValueOf(oldMap[i]).Pointer() {
			return false
		}
	}

	return true
}

func compareHooks(newHooks, oldHooks *Component) bool {
	return compareHook(newHooks.Hooks.Created, oldHooks.Hooks.Created) &&
		compareHook(newHooks.Hooks.Mounted, oldHooks.Hooks.Mounted) &&
		compareHookWithControl(newHooks.Hooks.BeforeCreated, oldHooks.Hooks.BeforeCreated) &&
		compareHook(newHooks.Hooks.BeforeDestroy, oldHooks.Hooks.BeforeDestroy) &&
		compareHook(newHooks.Hooks.BeforeUpdate, oldHooks.Hooks.BeforeUpdate) &&
		compareHook(newHooks.Hooks.Updated, oldHooks.Hooks.Updated)
}

func compareHookWithControl(newHook, oldHook HookWithControl) bool {
	if newHook == nil && oldHook == nil {
		return true
	}

	if newHook == nil || oldHook == nil {
		return false
	}

	return reflect.ValueOf(newHook).Pointer() == reflect.ValueOf(oldHook).Pointer()
}

func compareHook(newHook, oldHook Hook) bool {
	if newHook == nil && oldHook == nil {
		return true
	}

	if newHook == nil || oldHook == nil {
		return false
	}

	return reflect.ValueOf(newHook).Pointer() == reflect.ValueOf(oldHook).Pointer()
}

func compareForDirectives(newC, oldC *Component) bool {
	/*
		It's really bad way to fix bug with not-updated i, el in components Methods.
		We can only ForceUpdate methods, binds, e.t.c. and don't ForceUpdate 'body', but it will BE in the future...
	*/

	newIsItem, newI, newVal := newC.ForItemInfo()
	oldIsItem, oldI, oldVal := oldC.ForItemInfo()

	if !newIsItem && !oldIsItem {
		return true
	}

	if newIsItem != oldIsItem {
		return false
	}

	return newI == oldI && reflect.DeepEqual(newVal, oldVal)
}
