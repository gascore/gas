package store

import (
	"errors"
	"github.com/gascore/gas/std/store"
	"github.com/gascore/gas/web"
	webStore "github.com/gascore/gas/web/store"
)

var localStorage = webStore.NewDataStore(webStore.JSONEncoding, web.GetLocalStore)

var S *store.Store

func InitStore() error {
	var err error

	handlers := store.NewHandlers()
	handlers.Add("updateCount", func(s *store.Store, values ...interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"count": s.Get("count").(int) + 1,
		}, nil
	})
	handlers.AddMany(listsHandlers)

	S, err = store.New(&store.Store{
		Data: map[string]interface{}{
			"count":     0,
			"all":       []interface{}{},
			"completed": []interface{}{},
			"deleted":   []interface{}{},
		},
		Handlers: handlers,
		OnCreate: []store.OnCreateHook{
			func(s *store.Store) error {
				var dataRaw interface{}
				err := localStorage.Get("data", &dataRaw)
				if err != nil && err != webStore.ErrNilValue {
					return err
				}

				if dataRaw == nil {
					return nil
				}

				data, ok := dataRaw.(map[string]interface{})
				if !ok {
					return errors.New("invalid data type")
				}

				s.Data = data
				return nil
			},
		},
		AfterEmit: []store.AfterEmitHook{
			func(s *store.Store, eventName string, updatesMap map[string]interface{}, values []interface{}) error {
				err := localStorage.Set("data", s.Data)
				if err != nil {
					return err
				}

				return nil
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
