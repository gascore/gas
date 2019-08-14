package gas

import (
	"fmt"
	"reflect"
)

// Changed return isChanged? and canGoDeeper?
func Changed(newEl, oldEl interface{}) (bool, bool, error) {
	if reflect.TypeOf(newEl) != reflect.TypeOf(oldEl) {
		return true, false, nil
	}

	switch newEl.(type) {
	case *Component:
		isEquals, canGoDeeper := isComponentsEquals(I2C(newEl), I2C(oldEl))
		return !isEquals, canGoDeeper, nil
	case *Element:
		isEquals, canGoDeeper := isNodesEquals(I2E(newEl), I2E(oldEl))
		return !isEquals, canGoDeeper, nil
	case bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return newEl != oldEl, false, nil
	case fmt.Stringer:
		return newEl.(fmt.Stringer).String() != oldEl.(fmt.Stringer).String(), false, nil
	default:
		return false, false, fmt.Errorf("changed: invalid `newEl` or `oldEl`. types: %T, %T", newEl, oldEl)
	}
}

func isComponentsEquals(newC, oldC *C) (bool, bool) {
	isEquals := newC.ElementIsImportant == oldC.ElementIsImportant &&
		newC.RefsAllowed == oldC.RefsAllowed &&
		compareHooks(newC.Hooks, oldC.Hooks)

	if isEquals && newC.ElementIsImportant {
		return isElementsEquals(newC.Element, oldC.Element)
	}

	return isEquals, isEquals
}

func isNodesEquals(newE, oldE *E) (bool, bool) {
	if newE.Component != nil || oldE.Component != nil {
		if oldE.Component == nil || newE.Component == nil {
			return false, false
		}

		return isComponentsEquals(newE.Component, oldE.Component)
	}

	return isElementsEquals(newE, oldE)
}

func isElementsEquals(newE, oldE *E) (bool, bool) {
	canBeUpdated := newE.Tag == oldE.Tag && newE.HTML.Rendered == oldE.HTML.Rendered && reflect.ValueOf(newE.HTML.Render).Pointer() == reflect.ValueOf(oldE.HTML.Render).Pointer()
	return canBeUpdated && compareAttributes(newE.RAttrs, oldE.RAttrs), canBeUpdated
}

func compareAttributes(newMap, oldMap Map) bool {
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

func compareHooks(newHooks, oldHooks Hooks) bool {
	return compareHook(newHooks.Created, oldHooks.Created) &&
		compareHook(newHooks.Mounted, oldHooks.Mounted) &&
		compareHookWithControl(newHooks.BeforeCreated, oldHooks.BeforeCreated) &&
		compareHook(newHooks.BeforeDestroy, oldHooks.BeforeDestroy) &&
		compareHook(newHooks.BeforeUpdate, oldHooks.BeforeUpdate) &&
		compareHook(newHooks.Updated, oldHooks.Updated)
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

func DiffAttrs(newA, oldA Map) Map {
	diffMap := make(Map)
	for key, val := range oldA {
		if _, ok := diffMap[key]; ok {
			continue
		}

		if newA[key] != val {
			if _, ok := newA[key]; ok {
				diffMap[key] = newA[key]
				continue
			}

			diffMap[key] = ""
		}
	}

	for key, val := range newA {
		if _, ok := diffMap[key]; ok {
			continue
		}

		if oldA[key] != val {
			diffMap[key] = val
		}
	}
	return diffMap
}
