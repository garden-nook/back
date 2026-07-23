package repository

import (
	"context"
	"garden-nook/internal/modules/plot/dto"
	"garden-nook/internal/modules/plot/model"
	"garden-nook/internal/pkg/database"

	"github.com/jackc/pgx/v5"
)

type SnapshotRepo struct {
	db     database.DBTX
	mapper *database.ErrorMapper
}

func NewSnapshotRepo(db database.DBTX, mapper *database.ErrorMapper) *SnapshotRepo {
	return &SnapshotRepo{db: db, mapper: mapper}
}

func (r *SnapshotRepo) WithTx(tx pgx.Tx) *SnapshotRepo {
	return &SnapshotRepo{db: tx, mapper: r.mapper}
}

// CreateSnapshot вставляет новый снимок и возвращает его ID.
func (r *SnapshotRepo) CreateSnapshot(ctx context.Context, req dto.CreateSnapshotRequest) (string, error) {
	var snapshotID string
	err := r.db.QueryRow(ctx,
		`INSERT INTO snapshots (plot_id, snapshot_date, snapshot_type, grid_state, beds_state, objects_state, last_event_sequence)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING snapshot_id`,
		req.PlotID, req.SnapshotDate, req.SnapshotType,
		req.GridState, req.BedsState, req.ObjectsState,
		req.LastEventSequence,
	).Scan(&snapshotID)
	if err != nil {
		return "", r.mapper.Map(err)
	}
	return snapshotID, nil
}

// GetSnapshotByID возвращает снимок по идентификатору.
func (r *SnapshotRepo) GetSnapshotByID(ctx context.Context, snapshotID string) (*model.Snapshot, error) {
	row, err := r.db.Query(ctx,
		`SELECT snapshot_id, plot_id, snapshot_date, snapshot_type, grid_state, beds_state, objects_state, last_event_sequence, created_at
		 FROM snapshots WHERE snapshot_id = $1`, snapshotID)
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	defer row.Close()

	snap, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[model.Snapshot])
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	return snap, nil
}

// GetLatestSnapshot возвращает последний (по дате и, возможно, по sequence_number) снимок для участка.
func (r *SnapshotRepo) GetLatestSnapshot(ctx context.Context, plotID string) (*model.Snapshot, error) {
	row, err := r.db.Query(ctx,
		`SELECT snapshot_id, plot_id, snapshot_date, snapshot_type, grid_state, beds_state, objects_state, last_event_sequence, created_at
		 FROM snapshots
		 WHERE plot_id = $1
		 ORDER BY snapshot_date DESC, created_at DESC
		 LIMIT 1`, plotID)
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	defer row.Close()

	snap, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[model.Snapshot])
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	return snap, nil
}

// GetSnapshotsByPlot возвращает список снимков участка с пагинацией (сортировка по дате убыванию).
func (r *SnapshotRepo) GetSnapshotsByPlot(ctx context.Context, plotID string, p *database.Pagination) ([]model.Snapshot, int, error) {
	base := `SELECT snapshot_id, plot_id, snapshot_date, snapshot_type, grid_state, beds_state, objects_state, last_event_sequence, created_at
		FROM snapshots WHERE plot_id = $1`
	filterArgs := []any{plotID}
	order := ` ORDER BY snapshot_date DESC, created_at DESC`

	if p == nil {
		rows, err := r.db.Query(ctx, base+order, filterArgs...)
		if err != nil {
			return nil, 0, r.mapper.Map(err)
		}
		defer rows.Close()
		snaps, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Snapshot])
		if err != nil {
			return nil, 0, r.mapper.Map(err)
		}
		return snaps, len(snaps), nil
	}

	batch := &pgx.Batch{}
	countSQL := `SELECT COUNT(*) FROM snapshots WHERE plot_id = $1`
	batch.Queue(countSQL, plotID)

	pagSQL, pagArgs := p.SQL(2) // $2, $3 для limit/offset
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

	snaps, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Snapshot])
	if err != nil {
		return nil, 0, r.mapper.Map(err)
	}
	return snaps, total, nil
}
