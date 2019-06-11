package store

import (
	"errors"
	"fmt"
	"github.com/gascore/gas"
	"reflect"
	"strings"
)

// Store main structure
type Store struct {
	Data     map[string]interface{}
	Handlers map[string]Handler

	MiddleWares []MiddleWare

	OnCreate   []OnCreateHook
	BeforeEmit []BeforeEmitHook
	AfterEmit  []AfterEmitHook

	subscribers []Sub

	forRootParents    map[*gas.Component]bool
	forNonRootParents map[*gas.Component]bool
}

// MiddleWare let you do something before all events who have this (MiddleWare.Prefix) prefix.
//
// Example: { Prefix: "hello", Hook: func(s *Store) error { log.Println("Someone said hello }.
// This middleware will trigger on events: "helloMark", "helloElen", "helloArtem", "hello*etc*"
type MiddleWare struct {
	Prefix string
	Hook   func(s *Store, values []interface{}) error
}

// Hooks - functions dispatching by gas-store in special moments
type OnCreateHook func(s *Store) error
type BeforeEmitHook func(store *Store, eventName string, values []interface{}) error
type AfterEmitHook func(store *Store, eventName string, updatesMap map[string]interface{}, values []interface{}) error

// Handler - event handler with your stuff. Returns updatesData which will be appended to main store Data.
type Handler func(s *Store, values ...interface{}) (updatesMap map[string]interface{}, err error)

// Sub - component which will call ForceUpdate after Store updates
type Sub struct {
	Component *gas.Component

	CustomUpdater func() bool
}

// New initialize new store
func New(s *Store) (*Store, error) {
	if s.OnCreate != nil {
		for _, create := range s.OnCreate {
			err := create(s)
			if err != nil {
				return nil, err
			}
		}
	}

	if s.Data == nil {
		return nil, errors.New("store data is nil")
	}

	s.forRootParents = make(map[*gas.Component]bool)
	s.forNonRootParents = make(map[*gas.Component]bool)

	return s, nil
}

// GetSafely return Store.Data value by query
func (s *Store) GetSafely(query string) (interface{}, error) {
	val := s.Data[query]
	if val == nil {
		return nil, fmt.Errorf("undefined value: %s", query)
	}

	return val, nil
}

// Get proxy for GetSafely with error ignoring
func (s *Store) Get(query string) interface{} {
	val, _ := s.GetSafely(query)
	return val
}

// Emit runs event from Store handlers
func (s *Store) Emit(query string, values ...interface{}) error {
	handler := s.Handlers[query]
	if handler == nil {
		return fmt.Errorf("undefined event name: %s", query)
	}

	if s.BeforeEmit != nil {
		for _, beforeEmit := range s.BeforeEmit {
			if err := beforeEmit(s, query, values); err != nil {
				return nil
			}
		}
	}

	for _, mw := range s.MiddleWares {
		if !strings.HasPrefix(query, mw.Prefix) {
			continue
		}

		if mw.Hook == nil {
			return fmt.Errorf("hook is nil in middleware with prefix '%s'", mw.Prefix)
		}

		err := mw.Hook(s, values)
		if err != nil {
			return err
		}
	}

	updatesMap, err := handler(s, values...)
	if err != nil {
		return err
	}

	if updatesMap == nil {
		return nil
	}

	err = s.UpdateStore(updatesMap)
	if err != nil {
		return err
	}

	if s.AfterEmit != nil {
		for _, afterEmit := range s.AfterEmit {
			if err := afterEmit(s, query, updatesMap, values); err != nil {
				return nil
			}
		}
	}

	return nil
}

// UpdateStore update Store by replacing fields from updatesMap to Store.data
func (s *Store) UpdateStore(updatesMap map[string]interface{}) error {
	for uKey, uValue := range updatesMap {
		oValue := s.Data[uKey]
		if oValue == nil {
			return fmt.Errorf("undefined field in Data: %s", uKey)
		}

		if reflect.TypeOf(uValue) != reflect.TypeOf(oValue) {
			return fmt.Errorf("uncompared fields: %T and %T", uValue, oValue)
		}

		s.Data[uKey] = uValue
	}

	return s.update()
}

// RegisterComponent register new component in store
func (s *Store) RegisterComponent(c *gas.Component, customUpdater func() bool) *gas.Component {
	sub := Sub{
		Component:     c,
		CustomUpdater: customUpdater,
	}

	mounted := c.Hooks.Mounted
	c.Hooks.Mounted = func(this *gas.Component) error {
		if s.isRoot(sub.Component) {
			s.subscribers = append(s.subscribers, sub)
		}

		if mounted != nil {
			err := mounted(this)
			if err != nil {
				return err
			}
		}

		return nil
	}

	willDestroy := c.Hooks.BeforeDestroy
	c.Hooks.BeforeDestroy = func(this *gas.Component) error {
		for i, c := range s.subscribers {
			if sub.Component == c.Component {
				// remove sub from subscribers
				s.subscribers = append(s.subscribers[0:i], s.subscribers[i+1:]...)
			}
		}

		if willDestroy != nil {
			err := willDestroy(this)
			if err != nil {
				return err
			}
		}

		return nil
	}

	return c
}

// RC alias for Store.RegisterComponent
func (s *Store) RC(c *gas.Component, customUpdater func() bool) *gas.Component {
	return s.RegisterComponent(c, customUpdater)
}

// isRoot check if component have no RegisteredComponents which will update him after store updates
func (s *Store) isRoot(c *gas.Component) bool {
	parent := findParent(c)

	if parent == nil {
		s.forRootParents[parent] = true
		return true
	}

	if s.forRootParents[parent] {
		return true
	}
	if s.forNonRootParents[parent] {
		return false
	}

	var haveParentIsSubs bool
	for _, sub := range s.subscribers {
		if sub.Component == parent {
			haveParentIsSubs = true
		}
	}

	if haveParentIsSubs {
		s.forNonRootParents[parent] = true
		return false
	} else {
		return s.isRoot(parent)
	}
}

// findParent find component *true* parent which is not Element
func findParent(c *gas.Component) *gas.Component {
	if c.Parent == nil {
		return c.Parent
	}

	if c.Parent.IsElement() {
		return findParent(c.Parent)
	}

	return c.Parent
}

// update run ForceUpdate for all subs
func (s *Store) update() error {
	for _, sub := range s.subscribers {
		if sub.CustomUpdater != nil && !sub.CustomUpdater() {
			return nil
		}

		if sub.Component.GetElementUnsafely() == nil {
			return errors.New("element undefined")
		}

		err := sub.Component.ForceUpdate()
		if err != nil {
			return err
		}
	}

	return nil
}
