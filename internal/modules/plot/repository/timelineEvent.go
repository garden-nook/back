package repository

import (
	"context"
	"garden-nook/internal/modules/plot/dto"
	"garden-nook/internal/modules/plot/model"
	"garden-nook/internal/pkg/database"

	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

type TimelineRepo struct {
	db     database.DBTX
	mapper *database.ErrorMapper
}

func NewTimelineRepo(db database.DBTX, mapper *database.ErrorMapper) *TimelineRepo {
	return &TimelineRepo{db: db, mapper: mapper}
}

func (r *TimelineRepo) WithTx(tx pgx.Tx) *TimelineRepo {
	return &TimelineRepo{db: tx, mapper: r.mapper}
}

func (r *TimelineRepo) AddTimelineEvent(ctx context.Context, req dto.AddTimelineEventRequest) (string, error) {
	var eventID string
	err := r.db.QueryRow(ctx,
		`INSERT INTO timeline_events 
		 (plot_id, event_type, occurred_at, display_title, display_description, affected_cells, metadata, source_event_ids)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING event_id`,
		req.PlotID, req.EventType, req.OccurredAt,
		req.DisplayTitle, req.DisplayDescription,
		pq.Array(req.AffectedCells),
		req.Metadata,
		pq.Array(req.SourceEventIDs),
	).Scan(&eventID)
	if err != nil {
		return "", r.mapper.Map(err)
	}
	return eventID, nil
}

// GetTimelineEvents возвращает события ленты для участка с пагинацией, сортировка по occurred_at DESC.
func (r *TimelineRepo) GetTimelineEvents(ctx context.Context, plotID string, p *database.Pagination) ([]model.TimelineEvent, int, error) {
	base := `SELECT event_id, plot_id, event_type, occurred_at, display_title, display_description, affected_cells, metadata, source_event_ids
		FROM timeline_events WHERE plot_id = $1`
	filterArgs := []any{plotID}
	order := ` ORDER BY occurred_at DESC`

	if p == nil {
		rows, err := r.db.Query(ctx, base+order, filterArgs...)
		if err != nil {
			return nil, 0, r.mapper.Map(err)
		}
		defer rows.Close()
		events, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.TimelineEvent])
		if err != nil {
			return nil, 0, r.mapper.Map(err)
		}
		return events, len(events), nil
	}

	batch := &pgx.Batch{}
	countSQL := `SELECT COUNT(*) FROM timeline_events WHERE plot_id = $1`
	batch.Queue(countSQL, plotID)

	pagSQL, pagArgs := p.SQL(2)
	dataSQL := base + order + pagSQL
	dataArgs := append(filterArgs, pagArgs...)
	batch.Queue(dataSQL, dataArgs...)

	results := r.db.SendBatch(ctx, batch)
	defer results.Close()

	var total int
	if err := results.QueryRow().Scan(&total); err != nil {
		return nil, 0, r.mapper.Map(err)
	}
	rows, err := results.Query()
	if err != nil {
		return nil, 0, r.mapper.Map(err)
	}
	defer rows.Close()

	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.TimelineEvent])
	if err != nil {
		return nil, 0, r.mapper.Map(err)
	}
	return events, total, nil
}
