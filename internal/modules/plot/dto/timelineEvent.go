package dto

import (
	"encoding/json"
	"time"
)

type AddTimelineEventRequest struct {
	PlotID             string
	EventType          int
	OccurredAt         time.Time
	DisplayTitle       string
	DisplayDescription *string
	AffectedCells      []string
	Metadata           json.RawMessage
	SourceEventIDs     []string
}
