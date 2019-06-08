package dnd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gascore/dom"
	"github.com/gascore/gas"
	web "github.com/gascore/gas/web"
	sjs "syscall/js"
	"github.com/pkg/errors"
)

// Config dnd-lists config
type Config struct {
	This *gas.C // link to Compoent storing arrays
	Tag  string

	Group     string
	FieldName string

	Events EventsHandlers

	GroupClass       string
	PreviewClass     string // class for the dragging item preview
	PlaceholderClass string // class for the dragging item placeholder
}

// Lists return dnd-lists component
func Lists(config Config, e gas.External) *gas.C {
	switch {
	case config.This == nil:
		dom.ConsoleError("invalid This: config.This == nil")
		return nil
	case len(e.Body) == 0:
		dom.ConsoleError("no items in DragAndDrop body")
		return nil
	case len(config.Group) == 0:
		dom.ConsoleError("invalid Group name")
		return nil
	case len(config.FieldName) == 0:
		dom.ConsoleError("invalid FieldName name")
		return nil
	case len(config.PlaceholderClass) == 0:
		config.PlaceholderClass = "dnd-placeholder"
	case len(config.PreviewClass) == 0:
		config.PreviewClass = "dnd-preview"
	}

	body := gas.UnSpliceBody(e.Body)
	for i, item := range body {
		if !gas.IsComponent(item) {
			dom.ConsoleError("invalid body child type (not *gas.Component)")
			return nil
		}
		item := gas.I2C(item)

		if item.Handlers == nil {
			item.Handlers = make(map[string]gas.Handler)
		}
		if item.Attrs == nil {
			item.Attrs = make(map[string]string)
		}

		item.Attrs["draggable"] = "true"
		item.Attrs["data-group"] = config.Group
		item.Attrs["data-field"] = config.FieldName
		item.Attrs["data-is-item"] = "true"
		item.Attrs["data-dnd-index"] = fmt.Sprintf("%d", i)

		item.Hooks.Mounted = func(this *gas.Component) error {
			_el := this.Element().(*dom.Element)

			_el.AddEventListener("dragstart", func(e dom.Event) {
				_el := e.Target()

				indexS := _el.GetAttribute("data-dnd-index").String()
				index, err := strconv.Atoi(indexS)
				if err != nil {
					dom.ConsoleError(err.Error())
					return
				}

				if config.Events.Started != nil {
					block, err := config.Events.Started(StartedEvent{
						Index: index,
						Body:  web.ToUniteObject(e),
					})
					if err != nil {
						dom.ConsoleError(err.Error())
						return
					}

					if block {
						e.PreventDefault()
						return
					}
				}

				dataTransfer := e.JSValue().Get("dataTransfer")
				dataTransfer.Call("setData", "group", config.Group)
				dataTransfer.Call("setData", "field", config.FieldName)
				dataTransfer.Call("setData", "index", indexS)
				dataTransfer.Call("setData", "uuid", _el.GetAttribute("data-i").String())
				dataTransfer.Call("setData", "group-id", _el.ParentElement().GetAttribute("data-dnd-group-id").String())
				dataTransfer.Set("effectAllowed", "move")

				_p := _el.ParentElement()
				for _, _x := range dom.Doc.QuerySelectorAll("." + config.GroupClass) {
					if _x.JSValue() == _p.JSValue() {
						continue
					}

					_x.ClassList().Add(config.GroupClass + "-receiving")
				}
				_p.ClassList().Add(config.GroupClass + "-dragging")

				_preview := _el.Clone()
				_preview.ClassList().Add(config.PreviewClass)
				_preview.Style().Set("position", "absolute")
				_preview.Style().Set("top", "0")
				_preview.Style().Set("left", "0")
				_preview.Style().Set("zIndex", "-1")

				_el.AppendChild(_preview)
				dataTransfer.Call("setDragImage", _preview.JSValue(), e.JSValue().Get("offsetX").Int()+10, e.JSValue().Get("offsetY").Int()+10)

				go func() {
					time.Sleep(50 * time.Millisecond)
					_preview.ParentElement().RemoveChild(_preview)
				}()

				e.Target().ClassList().Add(config.PlaceholderClass)
			})

			_el.AddEventListener("drop", func(e dom.Event) {
				err := dropEvent(_el, config, e)
				if err != nil {
					dom.ConsoleError(err.Error())
					return
				}
			})

			_el.AddEventListener("dragend", func(event dom.Event) {
				target := event.Target()

				target.ClassList().Remove(config.PlaceholderClass)

				indexS := target.GetAttribute("data-dnd-index").String()
				oldIndex, err := strconv.Atoi(indexS)
				if err != nil {
					this.WarnError(err)
					return
				}

				newIndex := func() int {
					_childes := this.Parent.Element().(*dom.Element).ChildNodes()
					elID := target.GetAttribute("data-i").String()
					for i, el := range _childes {
						if el.GetAttribute("data-i").String() == elID {
							return i
						}
					}
					return -1
				}()

				if newIndex == -1 {
					if config.Events.Removed != nil {
						this.WarnError(
							config.Events.Removed(
								RemovedEvent{
									OldIndex: oldIndex,
									Body:     web.ToUniteObject(event),
								},
							),
						)
					}
				}

				if config.Events.Moved != nil {
					this.WarnError(
						config.Events.Moved(
							StandardEvent{
								OldIndex: oldIndex,
								NewIndex: newIndex,
								Body:     web.ToUniteObject(event),
							},
						),
					)
				}
			})

			return nil
		}
	}

	return gas.NC(
		&gas.C{
			Tag: config.Tag,
			Hooks: gas.Hooks{
				Mounted: func(this *gas.Component) error {
					_el := this.Element().(*dom.Element)

					_el.AddEventListener("drop", func(e dom.Event) {
						err := dropEvent(_el, config, e)
						if err != nil {
							dom.ConsoleError(err.Error())
							return
						}
					})

					_el.AddEventListener("dragenter", func(e dom.Event) {
						_placeholder := getPlaceholderNode(config)
						if _placeholder == nil {
							dom.ConsoleError("placeholder not found")
							return
						}

						_target := e.Target()
						targetTag := strings.ToLower(_target.TagName())

						if targetTag == "ul" {
							if !_target.ClassList().Contains(config.GroupClass) {
								return
							}
						} else if _target.GetAttribute("data-is-item").String() != "true" || _target.JSValue() == _placeholder.JSValue() {
							return
						}

						if _placeholder.GetAttribute("data-group").String() != config.Group {
							return
						}

						if config.Events.Entered != nil {
							aField := _placeholder.GetAttribute("data-field").String()
							oldIndexS := _placeholder.GetAttribute("data-dnd-index").String()
							oldIndex, err := strconv.Atoi(oldIndexS)
							if err != nil {
								dom.ConsoleError(err.Error())
								return
							}

							block, err := config.Events.Entered(EnteredEvent{Index: oldIndex, FieldName: aField, Body: web.ToUniteObject(e)})
							if err != nil {
								dom.ConsoleError(err.Error())
								return
							}

							if block {
								return
							}
						}

						e.PreventDefault()

						if targetTag == "ul" { // list
							// append placeholder to current list
							_target.AppendChild(_placeholder)
						} else { // item
							// move placeholder to a new place
							_pcParent := _placeholder.ParentElement()

							var insertAfter bool
							if _placeholder.ParentElement().JSValue() == _target.ParentElement().JSValue() {
								placeholderIndex := elementIndex(_placeholder)
								if placeholderIndex == -1 {
									dom.ConsoleError("cannot get _placeholder index in parent")
									return
								}

								targetIndex := elementIndex(_target)
								if targetIndex == -1 {
									dom.ConsoleError("cannot get _target index in parent")
									return
								}

								if placeholderIndex == targetIndex-1 { // _target after _placeholder
									insertAfter = true
								}
							}

							_pcParent.ClassList().Remove(config.GroupClass + "-dragging")
							_pcParent.ClassList().Add(config.GroupClass + "-receiving")

							if insertAfter {
								_pcParent.InsertAfter(_placeholder, _target)
							} else {
								_target.ParentElement().InsertBefore(_placeholder, _target)
							}

							_target.ParentElement().ClassList().Remove(config.GroupClass + "-receiving")
							_target.ParentElement().ClassList().Add(config.GroupClass + "-dragging")
						}
					})

					_el.AddEventListener("dragleave", func(event dom.Event) {
						target := event.JSValue().Get("target")
						if target.Get("data").String() != "" {
							return
						}

						if config.Events.Leaved != nil {
							this.WarnError(
								config.Events.Leaved(
									LeavedEvent{
										Body: web.ToUniteObject(event),
									},
								),
							)
						}
					})

					_el.AddEventListener("dragover", func(event dom.Event) {
						event.PreventDefault()
						event.JSValue().Get("dataTransfer").Set("dropEffect", "move")
					})

					return nil
				},
			},
			Attrs: map[string]string{
				"data-dnd-field": config.FieldName,
				"class":          config.GroupClass,
			},
		},
		func(this *gas.Component) []interface{} {
			if e.Slots != nil && e.Slots["header"] != nil {
				body = append([]interface{}{e.Slots["header"]}, body...)
			}

			if e.Slots != nil && e.Slots["footer"] != nil {
				body = append(body, e.Slots["footer"])
			}

			return body
		})
}

func dropEvent(_x *dom.Element, config Config, event dom.Event) error {
	dataTransfer := event.JSValue().Get("dataTransfer")
	if dataTransfer.Type() == sjs.TypeNull {
		return nil
	}

	aField := dataTransfer.Call("getData", "field").String()
	oldIndexS := dataTransfer.Call("getData", "index").String()
	oldIndex, err := strconv.Atoi(oldIndexS)
	if err != nil {
		return err
	}

	_aGroup := dom.Doc.QuerySelector("[data-dnd-field='" + aField + "']")
	if _aGroup == nil {
		return errors.New("another group element undefined")
	}

	_placeholder := getPlaceholderNode(config)
	if _placeholder == nil {
		return errors.New("placeholder not found")
	}

	if _x.GetAttribute("data-i").String() == _placeholder.GetAttribute("data-i").String() {
		return nil
	}

	if dataTransfer.Call("getData", "group").String() != config.Group { //|| // if groups are not same
		return movePlaceholder(config, oldIndex, _aGroup, _placeholder)
	}

	if _placeholder.ParentElement().GetAttribute("data-dnd-field").String() != config.FieldName { // little bug created by dragend
		return movePlaceholder(config, oldIndex, _aGroup, _placeholder)
	}

	event.PreventDefault()
	event.StopPropagation()

	newIndex := elementIndex(_placeholder)
	if newIndex == -1 {
		return errors.New("_placeholder new place undefined")
	}

	err = movePlaceholder(config, oldIndex, _aGroup, _placeholder)
	if err != nil {
		return err
	}

	eData := StandardEvent{
		OldIndex: oldIndex,
		NewIndex: newIndex,

		Body: web.ToUniteObject(event),
	}

	data := config.This.Get(config.FieldName).([]interface{})
	if aField == config.FieldName {
		if newIndex == oldIndex {
			return nil
		}

		// replace in this array
		config.This.SetValue(config.FieldName, replaceInArr(data, newIndex, oldIndex))

		if config.Events.Updated != nil {
			err := config.Events.Updated(eData)
			if err != nil {
				return err
			}
		}
	} else {
		aData := config.This.Get(aField).([]interface{})

		config.This.Set(map[string]interface{}{
			config.FieldName: append(data[:newIndex], append([]interface{}{aData[oldIndex]}, data[newIndex:]...)...),
			aField:           append(aData[:oldIndex], aData[oldIndex+1:]...),
		})

		if config.Events.Added != nil {
			err := config.Events.Added(eData)
			if err != nil {
				return err
			}
		}
	}

	if config.Events.Ended != nil {
		err := config.Events.Ended(eData)
		if err != nil {
			return err
		}
	}

	return nil
}

func getPlaceholderNode(config Config) *dom.Element {
	return dom.GetDocument().QuerySelector("." + config.PlaceholderClass)
}

func addToArr(arr []interface{}, el interface{}, i int) []interface{} {
	if len(arr) == 0 {
		return []interface{}{el}
	}

	arr = append(arr, el)
	copy(arr[i+1:], arr[i:])
	arr[i] = el

	return arr
}

func replaceInArr(arr []interface{}, new, old int) []interface{} {
	x := arr[old]
	out := addToArr(append(arr[:old], arr[old+1:]...), x, new)
	return out
}

func elementIndex(_el *dom.Element) int {
	childes := _el.ParentElement().ChildNodes()
	for i, _node := range childes {
		if _node.JSValue() == _el.JSValue() {
			return i
		}
	}
	return -1
}

func movePlaceholder(config Config, oldIndex int, _aGroup, _placeholder *dom.Element) error {
	err := func() error {
		childNodes := _aGroup.ChildNodes()
		if len(childNodes) == oldIndex {
			_aGroup.AppendChild(_placeholder)
			return nil
		}

		for i, _child := range childNodes {
			if i == oldIndex {
				_aGroup.InsertBefore(_placeholder, _child)
				return nil
			}
		}

		return errors.New("cannot move placeholder to old place")
	}()
	if err != nil {
		return err
	}

	for _, _el := range dom.Doc.QuerySelectorAll("." + config.GroupClass + "-dragging") {
		_el.ClassList().Remove(config.GroupClass + "-dragging")
	}

	for _, _el := range dom.Doc.QuerySelectorAll("." + config.GroupClass + "-receiving") {
		_el.ClassList().Remove(config.GroupClass + "-receiving")
	}

	return nil
}
