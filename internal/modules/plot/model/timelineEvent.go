package model

import (
	"encoding/json"
	"time"
)

type TimelineEvent struct {
	EventID            string          `json:"event_id" db:"event_id"`
	PlotID             string          `json:"plot_id" db:"plot_id"`
	EventType          int             `json:"event_type" db:"event_type"`
	OccurredAt         time.Time       `json:"occurred_at" db:"occurred_at"`
	DisplayTitle       string          `json:"display_title" db:"display_title"`
	DisplayDescription string          `json:"display_description,omitempty" db:"display_description"`
	AffectedCells      []string        `json:"affected_cells,omitempty" db:"affected_cells"`
	Metadata           json.RawMessage `json:"metadata" db:"metadata"`
	SourceEventIDs     []string        `json:"source_event_ids" db:"source_event_ids"` // uuid[]
}
