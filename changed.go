package gas

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
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
		newC := I2C(new)
		oldC := I2C(old)

		if newC.Directives.HTML.Rendered != oldC.Directives.HTML.Rendered {
			return true, nil
		}

		return !isComponentsEquals(newC, oldC), nil // thank you god for the go-cmp
	}

	return false, fmt.Errorf("changed: invalid `new` or `old`. types: %T, %T", new, old)
}

func isComponentsEquals(new, old *Component) bool {
	// sometimes i'm sad that i chose strict-type pl
	daE := true // cmp.Equal(new.Data, old.Data)
	wE  := cmp.Equal(new.Watchers, old.Watchers)
	mE  := true // cmp.Equal(new.Methods, old.Methods)
	coE := true // cmp.Equal(new.Computeds, old.Computeds)

	hE := compareHooks(new.Hooks, old.Hooks)
	bE := compareBinds(new.RenderedBinds, old.RenderedBinds)

	diIfE := reflect.ValueOf(new.Directives.If).Pointer() == reflect.ValueOf(old.Directives.If).Pointer()
	diFE  := compareForDirectives(new, old)
	diME  := (new.Directives.Model.Data == old.Directives.Model.Data) && (new.Directives.Model.Component == old.Directives.Model.Component)
	diHE  := reflect.ValueOf(new.Directives.HTML.Render).Pointer() == reflect.ValueOf(old.Directives.HTML.Render).Pointer()
	diE   := diIfE && diFE && diME && diHE // Directives

	tE := new.Tag == old.Tag
	aE := cmp.Equal(new.Attrs, old.Attrs)

	return daE && wE && mE && coE && hE && bE && diE && tE && aE
}

func compareHooks(new, old Hooks) bool {
	created := cmp.Equal(new.Created, old.Created)
	beforeCreate := cmp.Equal(new.BeforeCreate, old.BeforeCreate)
	destroyed := cmp.Equal(new.Destroyed, old.Destroyed)

	return created && beforeCreate && destroyed
}

func compareForDirectives(new, old *Component) bool {
	/*
		It's really bad way to fix bug with not-updated i, el in components Methods.
		We can only update methods, binds, e.t.c. and don't update 'body', but it will be in the future...
	*/

	newIsItem, newI, newVal := new.ForItemInfo()
	oldIsItem, oldI, oldVal := old.ForItemInfo()

	if newIsItem != oldIsItem {
		return false
	} else {
		return newI == oldI && newVal == oldVal
	}
}

func compareBinds(new, old map[string]string) bool {
	if len(new) != len(old) {
		return false
	}

	for newKey, newValue := range new {
		if newValue != old[newKey] {
			return false
		}
	}

	return true
}
