package dnd

import (
	"github.com/gascore/gas"
)

// EventHandlers struct for unit all event handlers. see events.md
//
// Event types:
//
// | Name      | Timing                                             | Structure       |
// |-----------|----------------------------------------------------|-----------------|
// | started   | Element is chosen by user (return true for block)  | StartedEvent    |
// | entered   | on dragenter (return true for block entering)      | EnteredEvent    |
// | leaved    | on dragleave                                		| LeavedEvent     |
// | ended     | Element dragging ended                             | StandardEvent   |
// | added     | Element is dropped into the list from another list | StandardEvent   |
// | updated   | When dragging element changes position             | StandardEvent   |
// | removed   | Element is removed from the list into another list | RemovedEvent    |
// | moved     | When you move an item in the list or between lists | StandardEvent   |
type EventsHandlers struct {
	Started func(event StartedEvent) (bool, error)
	Entered func(event EnteredEvent) (bool, error)
	Leaved  func(event LeavedEvent) error
	Ended   func(event StandardEvent) error
	Added   func(event StandardEvent) error
	Updated func(event StandardEvent) error
	Removed func(event RemovedEvent) error
	Moved   func(event StandardEvent) error
}

type StartedEvent struct {
	Index int

	Body gas.Object
}

type RemovedEvent struct {
	OldIndex int

	Body gas.Object
}

type EnteredEvent struct {
	Index     int
	FieldName string

	Body gas.Object
}

type LeavedEvent struct {
	Body gas.Object
}

type StandardEvent struct {
	OldIndex int
	NewIndex int

	Body gas.Object
}
