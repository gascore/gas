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
	case *Element:
		return !isNodesEquals(I2E(newEl), I2E(oldEl)), nil
	case bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return newEl != oldEl, nil
	case fmt.Stringer:
		return newEl.(fmt.Stringer).String() != oldEl.(fmt.Stringer).String(), nil
	default:
		return false, fmt.Errorf("changed: invalid `newEl` or `oldEl`. types: %T, %T", newEl, oldEl)
	}
}

func isComponentsEquals(newC, oldC *C) bool {
	isEquals := newC.ElementIsImportant == oldC.ElementIsImportant &&
		newC.RefsAllowed == oldC.RefsAllowed &&
		compareHooks(newC.Hooks, oldC.Hooks)

	if isEquals && newC.ElementIsImportant {
		return isElementsEquals(newC.Element, oldC.Element)
	}

	return isEquals
}

func isNodesEquals(newE, oldE *E) bool {
	if newE.Component != nil || oldE.Component != nil {
		if oldE.Component == nil || newE.Component == nil {
			return false
		}

		return isComponentsEquals(newE.Component, oldE.Component)
	}

	return isElementsEquals(newE, oldE)
}

func isElementsEquals(newE, oldE *E) bool {
	return ElementsCanBeUpdated(newE, oldE) && compareAttributes(newE.RAttrs, oldE.RAttrs)
}

func ElementsCanBeUpdated(newE, oldE *E) bool {
	return newE.Tag == oldE.Tag && 
		newE.HTML.Rendered == oldE.HTML.Rendered && 
		reflect.ValueOf(newE.HTML.Render).Pointer() == reflect.ValueOf(oldE.HTML.Render).Pointer()
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

func DiffAttrs(newA, oldA map[string]string) map[string]string {
	diffMap := make(map[string]string)
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
