package model

import (
	"encoding/json"
	"time"
)

type Event struct {
	EventID        string          `db:"event_id"`
	PlotID         string          `db:"plot_id"`
	EventType      int             `db:"event_type"`
	Payload        json.RawMessage `db:"payload"`
	OccurredAt     time.Time       `db:"occurred_at"`
	SequenceNumber int64           `db:"sequence_number"`
}
