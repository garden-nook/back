package enum

// EventType
// @Description.markdown eventtypes
type EventType int32

const (
	EventTypePlotCreated      EventType = 1
	EventTypePlotUpdated      EventType = 2
	EventTypePlotResized      EventType = 10
	EventTypeBedCreated       EventType = 11
	EventTypeBedUpdated       EventType = 12
	EventTypeBedDeleted       EventType = 13
	EventTypeObjectCreated    EventType = 14
	EventTypeObjectUpdated    EventType = 15
	EventTypeObjectDeleted    EventType = 16
	EventTypeCropPlanted      EventType = 17
	EventTypeCropRemoved      EventType = 18
	EventTypeCellShadeUpdated EventType = 19
)

func IsValidEventType(t EventType) bool {
	switch t {
	case EventTypeBedCreated, EventTypeBedUpdated, EventTypeBedDeleted,
		EventTypeObjectCreated, EventTypeObjectUpdated, EventTypeObjectDeleted,
		EventTypeCropPlanted, EventTypeCropRemoved, EventTypeCellShadeUpdated,
		EventTypePlotResized:
		return true
	default:
		return false
	}
}
