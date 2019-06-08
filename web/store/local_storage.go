// Copyright 2018 Noskov Artem, Alex Browne. All rights reserved.
// Use of this source code is governed by the MIT
// license, which can be found in the LICENSE file.

package store

import (
	"errors"
	"fmt"
	"github.com/gascore/gas"
)

// ErrLocalStorageNotSupported is returned if localStorage is not supported.
var ErrLocalStorageNotSupported = errors.New("localStorage does not appear to be supported/enabled in this browser")
var ErrNilValue = errors.New("localStorage item is null")

// ItemNotFoundError is returned if an item with the given key does not exist in
// localStorage.
type ItemNotFoundError struct {
	msg string
}

// Error implements the error interface.
func (e ItemNotFoundError) Error() string {
	return e.msg
}

func newItemNotFoundError(format string, args ...interface{}) ItemNotFoundError {
	return ItemNotFoundError{
		msg: fmt.Sprintf(format, args...),
	}
}

// DetectStorage detects and (re)initializes the localStorage.
func DetectStorage(getLocalStorage func() gas.Object) (ok bool, localStorage gas.Object) {
	defer func() {
		if r := recover(); r != nil {
			localStorage = nil
			ok = false
		}
	}()

	localStorage = getLocalStorage()

	if localStorage == nil {
		return false, nil
	}

	// Cf. https://developer.mozilla.org/en-US/docs/Web/API/Web_Storage_API/Using_the_Web_Storage_API
	// https://gist.github.com/paulirish/5558557
	x := "__storage_test__"
	obj := localStorage.Get(x)
	if obj == nil {
		localStorage = nil
		return false, nil
	}

	localStorage.Call("removeItem", x)
	return true, localStorage
}

// SetItem saves the given item in localStorage under the given key.
func SetItem(localStorage gas.Object, key, item string) (err error) {
	if localStorage == nil {
		return ErrLocalStorageNotSupported
	}
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("could not use local storage: %v", r)
			}
		}
	}()

	localStorage.Call("setItem", key, item)
	return nil
}

// GetItem finds and returns the item identified by key. If there is no item in
// localStorage with the given key, GetItem will return an ItemNotFoundError.
func GetItem(localStorage gas.Object, key string) (s string, err error) {
	if localStorage == nil {
		return "", ErrLocalStorageNotSupported
	}
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("could not use local storage: %v", r)
			}
			s = ""
		}
	}()

	item := localStorage.Call("getItem", key)
	if item == nil {
		err = newItemNotFoundError(
			"Could not find an item with the given key: %s", key)
	} else {
		s = item.String()
	}
	return s, err
}

// Key finds and returns the key associated with the given item. If the item is
// not in localStorage, Key will return an ItemNotFoundError.
func Key(localStorage gas.Object, item string) (s string, err error) {
	if localStorage == nil {
		return "", ErrLocalStorageNotSupported
	}
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("could not use local storage: %v", r)
			}
			s = ""
		}
	}()

	key := localStorage.Call("key", item)
	if key == nil {
		err = newItemNotFoundError(
			"Could not find a key for the given item: %s", item)
	} else {
		s = key.String()
	}
	return s, err
}

// RemoveItem removes the item with the given key from localStorage.
func RemoveItem(localStorage gas.Object, key string) (err error) {
	if localStorage == nil {
		return ErrLocalStorageNotSupported
	}
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("could not use local storage: %v", r)
			}
		}
	}()

	localStorage.Call("removeItem", key)
	return nil
}

// Length returns the number of items currently in localStorage.
func Length(localStorage gas.Object) (l int, err error) {
	if localStorage == nil {
		return 0, ErrLocalStorageNotSupported
	}
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("could not use local storage: %v", r)
			}
			l = 0
		}
	}()

	length := localStorage.GetInt("length")
	return length, nil
}

// Clear removes all items from localStorage.
func Clear(localStorage gas.Object) (err error) {
	if localStorage == nil {
		return ErrLocalStorageNotSupported
	}
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("could not use local storage: %v", r)
			}
		}
	}()

	localStorage.Call("clear")
	return nil
}
