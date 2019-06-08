// Copyright 2018 Noskov Artem, Alex Browne. All rights reserved.
// Use of this source code is governed by the MIT
// license, which can be found in the LICENSE file.

package store

import (
	"errors"
	"github.com/gascore/gas"
)

// DataStore is an object with methods for storing and retrieving arbitrary
// go data structures in localStorage.
type DataStore struct {
	Encoding     EncoderDecoder
	LocalStorage gas.Object
}

// NewDataStore creates and returns a new DataStore with the given encoding.
// locstor.JSON and locstor.Binary are two encodings provided by default. You
// can also pass in a custom encoding.
//
// For getLocalStorage you can use wasm.GetLocalStorage() (github.com/gascore/gas/web/wasm)
// or gojs.GetLocalStorage() (github.com/gascore/gas/web/wasm)
func NewDataStore(encoding EncoderDecoder, getLocalStorage func() gas.Object) *DataStore {
	ok, localStorage := DetectStorage(getLocalStorage)
	if !ok {
		panic(errors.New("cannot create new data store"))
	}

	return &DataStore{
		Encoding:     encoding,
		LocalStorage: localStorage,
	}
}

// Set saves the given item under the given key in localStorage.
func (store DataStore) Set(key string, item interface{}) error {
	encodedItem, err := store.Encoding.Encode(item)
	if err != nil {
		return err
	}
	return SetItem(store.LocalStorage, key, string(encodedItem))
}

// Get finds the item with the given key in localStorage and scans it into
// holder. holder must be a pointer to some data structure which is capable of
// holding the item. In general holder should be the same type as the item that
// was passed to Save.
func (store DataStore) Get(key string, holder interface{}) error {
	encodedItem, err := GetItem(store.LocalStorage, key)
	if err != nil {
		return err
	}

	if len(encodedItem) == 0 {
		return ErrNilValue
	}

	return store.Encoding.Decode([]byte(encodedItem), holder)
}

// Delete deletes the item with the given key from localStorage.
func (store DataStore) Delete(key string) error {
	return RemoveItem(store.LocalStorage, key)
}
