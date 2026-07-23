package dto

import (
	"encoding/json"
	"time"
)

type CreateSnapshotRequest struct {
	PlotID            string
	SnapshotDate      time.Time
	SnapshotType      int
	GridState         json.RawMessage
	BedsState         json.RawMessage
	ObjectsState      json.RawMessage
	LastEventSequence int64
}
