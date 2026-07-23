package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"garden-nook/internal/modules/plot/enum"
	"garden-nook/internal/pkg/database"

	"github.com/jackc/pgx/v5"
)

type EventStoreRepo struct {
	db     database.DBTX
	mapper *database.ErrorMapper
}

func NewEventStoreRepo(db database.DBTX, mapper *database.ErrorMapper) *EventStoreRepo {
	return &EventStoreRepo{db: db, mapper: mapper}
}

func (r *EventStoreRepo) WithTx(tx pgx.Tx) *EventStoreRepo {
	return &EventStoreRepo{db: tx, mapper: r.mapper}
}

func (r *EventStoreRepo) AppendEvent(ctx context.Context, plotID string, eventType enum.EventType, payload json.RawMessage) (int64, error) {
	lockKey := fmt.Sprintf("event_store:%s", plotID)
	batch := &pgx.Batch{}

	batch.Queue("SELECT pg_advisory_xact_lock(hashtext($1))", lockKey)

	batch.Queue(`
        INSERT INTO event_store (plot_id, event_type, payload, sequence_number)
        SELECT $1, $2, $3, COALESCE(MAX(sequence_number), 0) + 1
        FROM event_store
        WHERE plot_id = $1
        RETURNING sequence_number
    `, plotID, eventType, payload)

	br := r.db.SendBatch(ctx, batch)
	defer br.Close()

	if _, err := br.Exec(); err != nil {
		return 0, r.mapper.Map(err)
	}

	var seq int64
	if err := br.QueryRow().Scan(&seq); err != nil {
		return 0, r.mapper.Map(err)
	}
	return seq, nil
}
