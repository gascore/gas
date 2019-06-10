package store

import (
	"errors"
	"fmt"
	"github.com/gascore/gas/std/store"
)

var listsHandlers = map[string]store.Handler{
	"clearList": func(s *store.Store, values ...interface{}) (map[string]interface{}, error) {
		name, ok := values[0].(string)
		if !ok {
			return nil, errors.New("invalid list name")
		}

		return map[string]interface{}{
			name: []interface{}{},
		}, nil
	},
	"appendToList": func(s *store.Store, values ...interface{}) (map[string]interface{}, error) {
		el, ok := values[1].(string)
		if !ok {
			return nil, errors.New("invalid element value")
		}

		name, ok := values[0].(string)
		if !ok {
			return nil, errors.New("invalid 'name' type")
		}

		listUnTyped, err := s.GetSafely(name)
		if err != nil {
			return nil, err
		}

		list, ok := listUnTyped.([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid '%s' type", name)
		}

		list = append(list, el)

		return map[string]interface{}{
			name: list,
		}, nil
	},
	"editAll": func(s *store.Store, values ...interface{}) (map[string]interface{}, error) {
		i := values[0].(int)
		val := values[1].(string)

		all := s.Get("all").([]interface{})

		all[i] = val

		return map[string]interface{}{
			"all": all,
		}, nil
	},
	"deleteFromAll": func(s *store.Store, values ...interface{}) (map[string]interface{}, error) {
		i, ok := values[0].(int)
		if !ok {
			return nil, errors.New("invalid index")
		}

		appendToDeleted, ok := values[1].(bool)
		if !ok {
			return nil, errors.New("invalid appendToDeleted")
		}

		allUnTyped, err := s.GetSafely("all")
		if err != nil {
			return nil, err
		}

		all, ok := allUnTyped.([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid '%s' type", "allUnTyped")
		}

		el := all[i]

		copy(all[i:], all[i+1:]) // Shift a[i+1:] left one index
		all[len(all)-1] = ""     // Erase last element (write zero value)
		all = all[:len(all)-1]   // Truncate slice

		out := map[string]interface{}{
			"all": all,
		}

		if appendToDeleted {
			out["deleted"] = append(s.Get("deleted").([]interface{}), el)
		}

		return out, nil
	},
	"completeInAll": func(s *store.Store, values ...interface{}) (map[string]interface{}, error) {
		i, ok := values[0].(int)
		if !ok {
			return nil, errors.New("invalid index")
		}

		allUnTyped, err := s.GetSafely("all")
		if err != nil {
			return nil, err
		}

		all, ok := allUnTyped.([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid '%s' type", "allUnTyped")
		}

		el := all[i]

		copy(all[i:], all[i+1:]) // Shift a[i+1:] left one index
		all[len(all)-1] = ""     // Erase last element (write zero value)
		all = all[:len(all)-1]   // Truncate slice

		out := map[string]interface{}{
			"all": all,
		}

		out["completed"] = append(s.Get("completed").([]interface{}), el)

		return out, nil
	},
}
