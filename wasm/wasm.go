package wasm

import (
	"github.com/Sinicablyat/dom"
	"github.com/Sinicablyat/gas/core"
	"github.com/pkg/errors"
)

type BackEnd struct {}

func (w BackEnd) New(startPoint string) (string, error) {
	_el := dom.Doc.GetElementById(startPoint)
	if _el == nil {
		return "", errors.New("invalid startPoint")
	}

	return _el.GetTagName(), nil
}

func (w BackEnd) Init(gas core.Gas) error {
	app := gas.App
	_main := gas.GetElement().(*dom.Element)

	_main.SetAttribute("data-i", app.UUID)

	err := w.UpdateComponentChildes(&app, app.Childes(&app), []interface{}{})
	if err != nil {
		return err
	}

	dom.Doc.GetElementsByTagName("body")[0].SetAttribute("data-ready", true)

	return nil
}

func GetBackEnd() core.BackEnd {
	be := BackEnd{}
	return be
}
