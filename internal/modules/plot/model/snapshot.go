package model

import (
	"encoding/json"
	"time"
)

type Snapshot struct {
	SnapshotID        string          `db:"snapshot_id"`
	PlotID            string          `db:"plot_id"`
	SnapshotDate      time.Time       `db:"snapshot_date"`
	SnapshotType      int             `db:"snapshot_type"`
	GridState         json.RawMessage `db:"grid_state"`
	BedsState         json.RawMessage `db:"beds_state"`
	ObjectsState      json.RawMessage `db:"objects_state"`
	LastEventSequence int64           `db:"last_event_sequence"`
	CreatedAt         time.Time       `db:"created_at"`
}
