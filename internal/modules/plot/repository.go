package plots

import (
	"context"
	"fmt"
	"garden-nook/internal/pkg/apperrors"
	"garden-nook/internal/pkg/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db     database.DBTX
	mapper *database.ErrorMapper
}

func NewRepository(pool *pgxpool.Pool, mapper *database.ErrorMapper) *Repository {
	return &Repository{db: pool, mapper: mapper}
}

// WithTx возвращает копию репозитория, работающую в рамках переданной транзакции.
func (r *Repository) WithTx(tx pgx.Tx) *Repository {
	return &Repository{db: tx, mapper: r.mapper}
}

// =====================================================================
// PLOTS
// =====================================================================

func (r *Repository) ListPlots(ctx context.Context, ownerID string, p *database.Pagination) ([]Plot, int, error) {
	where := `WHERE p.owner_id = $1 AND p.is_deleted = FALSE`
	filterArgs := []any{ownerID}
	argIdx := 2

	baseSelect := `SELECT p.plot_id, p.name, p.soil_type, st.name AS soil_name,
	       			   p.area_sq_m, p.grid_cell_size, p.grid_cols, p.grid_rows
	                FROM plots p
	                JOIN soil_types st ON st.id = p.soil_type `

	if p == nil {
		rows, err := r.db.Query(ctx, baseSelect+where+` ORDER BY p.name`)
		if err != nil {
			return nil, 0, r.mapper.Map(err)
		}
		defer rows.Close()
		plots, err := pgx.CollectRows(rows, pgx.RowToStructByName[Plot])
		if err != nil {
			return nil, 0, r.mapper.Map(err)
		}
		return plots, len(plots), nil
	}

	batch := &pgx.Batch{}
	countSQL := `SELECT COUNT(*) FROM plots p ` + where
	batch.Queue(countSQL, filterArgs...)

	pagSQL, pagArgs := p.SQL(argIdx)
	dataSQL := baseSelect + where + ` ORDER BY p.name` + pagSQL
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
	plots, err := pgx.CollectRows(rows, pgx.RowToStructByName[Plot])
	if err != nil {
		return nil, 0, r.mapper.Map(err)
	}
	return plots, total, nil
}

func (r *Repository) CreatePlot(ctx context.Context, req *CreatePlotModel, ownerID string) (string, error) {
	var plotID string
	err := r.db.QueryRow(ctx,
		`INSERT INTO plots (name, owner_id, soil_type, boundary, area_sq_m,
			            grid_cell_size, grid_cols, grid_rows)
		 VALUES ($1, $2, $3,
		         ST_MakeEnvelope(0, 0, $4, $5, 3857),
		         $6, $7, $8, $9)
		 RETURNING plot_id`,
		req.Name, ownerID, req.SoilTypeID,
		req.BoundaryWidth, req.BoundaryHeight,
		req.AreaSqM, req.GridCellSize, req.GridCols, req.GridRows,
	).Scan(&plotID)
	if err != nil {
		return "", r.mapper.Map(err)
	}
	return plotID, nil
}

func (r *Repository) UpdatePlot(ctx context.Context, plotID, ownerID string, req UpdatePlotRequest) (string, error) {
	fields := []database.SetField{
		{Name: "name", Value: req.Name},
		{Name: "soil_type", Value: req.SoilTypeID},
	}
	setSQL, setArgs := database.BuildUpdateSet(1, fields...)
	if len(setArgs) == 0 {
		return plotID, nil
	}
	query := fmt.Sprintf(
		"UPDATE plots SET %s WHERE plot_id = $%d AND owner_id = $%d AND is_deleted = FALSE RETURNING plot_id",
		setSQL, len(setArgs)+1, len(setArgs)+2,
	)
	args := append(setArgs, plotID, ownerID)

	var updatedID string
	err := r.db.QueryRow(ctx, query, args...).Scan(&updatedID)
	if err != nil {
		return "", r.mapper.Map(err)
	}
	return updatedID, nil
}

func (r *Repository) SoftDeletePlot(ctx context.Context, plotID, ownerID string) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE plots SET is_deleted = TRUE WHERE plot_id = $1 AND owner_id = $2 AND is_deleted = FALSE`,
		plotID, ownerID)
	if err != nil {
		return r.mapper.Map(err)
	}
	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

//// =====================================================================
//// BEDS
//// =====================================================================
//
//func (r *Repository) CreateBed(ctx context.Context, b *Bed) error {
//	return r.db.QueryRow(ctx,
//		`INSERT INTO beds_ui (bed_id, plot_id, name, geom)
//		 VALUES ($1, $2, $3, ST_GeomFromGeoJSON($4))
//		 RETURNING created_at, updated_at`,
//		b.BedID, b.PlotID, b.Name, b.Geom,
//	).Scan(&b.CreatedAt, &b.UpdatedAt)
//}
//
//func (r *Repository) GetBedByID(ctx context.Context, bedID, plotID string) (*Bed, error) {
//	b := &Bed{}
//	var geom string
//	err := r.db.QueryRow(ctx,
//		`SELECT bed_id, plot_id, name, ST_AsGeoJSON(geom), is_deleted, created_at, updated_at
//		 FROM beds_ui WHERE bed_id=$1 AND plot_id=$2 AND is_deleted=FALSE`,
//		bedID, plotID,
//	).Scan(&b.BedID, &b.PlotID, &b.Name, &geom, &b.IsDeleted, &b.CreatedAt, &b.UpdatedAt)
//	if errors.Is(err, pgx.ErrNoRows) {
//		return nil, apperrors.ErrNotFound
//	}
//	if err != nil {
//		return nil, err
//	}
//	b.Geom = geom
//	return b, nil
//}
//
//func (r *Repository) UpdateBed(ctx context.Context, bedID, plotID string, name *string, geom *string) error {
//	setClauses := []string{}
//	args := []interface{}{}
//	argIdx := 1
//
//	if name != nil {
//		setClauses = append(setClauses, fmt.Sprintf("name=$%d", argIdx))
//		args = append(args, *name)
//		argIdx++
//	}
//	if geom != nil {
//		setClauses = append(setClauses, fmt.Sprintf("geom=ST_GeomFromGeoJSON($%d)", argIdx))
//		args = append(args, *geom)
//		argIdx++
//	}
//
//	if len(setClauses) == 0 {
//		return nil
//	}
//
//	setClauses = append(setClauses, "updated_at=NOW()")
//	query := fmt.Sprintf(
//		`UPDATE beds_ui SET %s WHERE bed_id=$%d AND plot_id=$%d AND is_deleted=FALSE`,
//		joinStrings(setClauses, ", "), argIdx, argIdx+1,
//	)
//	args = append(args, bedID, plotID)
//
//	tag, err := r.db.Exec(ctx, query, args...)
//	if err != nil {
//		return err
//	}
//	if tag.RowsAffected() == 0 {
//		return apperrors.ErrNotFound
//	}
//	return nil
//}
//
//func (r *Repository) SoftDeleteBed(ctx context.Context, bedID, plotID string) error {
//	tag, err := r.db.Exec(ctx,
//		`UPDATE beds_ui SET is_deleted=TRUE, updated_at=NOW()
//		 WHERE bed_id=$1 AND plot_id=$2 AND is_deleted=FALSE`,
//		bedID, plotID)
//	if err != nil {
//		return err
//	}
//	if tag.RowsAffected() == 0 {
//		return apperrors.ErrNotFound
//	}
//	return nil
//}
//
//func (r *Repository) ListBedsByPlot(ctx context.Context, plotID string) ([]Bed, error) {
//	rows, err := r.db.Query(ctx,
//		`SELECT bed_id, plot_id, name, ST_AsGeoJSON(geom), created_at, updated_at
//		 FROM beds_ui WHERE plot_id=$1 AND is_deleted=FALSE
//		 ORDER BY created_at`,
//		plotID)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var beds []Bed
//	for rows.Next() {
//		var b Bed
//		var geom string
//		if err := rows.Scan(&b.BedID, &b.PlotID, &b.Name, &geom, &b.CreatedAt, &b.UpdatedAt); err != nil {
//			return nil, err
//		}
//		b.Geom = geom
//		beds = append(beds, b)
//	}
//	return beds, rows.Err()
//}
//
//// =====================================================================
//// OBJECTS
//// =====================================================================
//
//func (r *Repository) CreateObject(ctx context.Context, o *UIObject) error {
//	return r.db.QueryRow(ctx,
//		`INSERT INTO objects_ui (object_id, plot_id, name, object_type, geom)
//		 VALUES ($1, $2, $3, $4, ST_GeomFromGeoJSON($5))
//		 RETURNING created_at, updated_at`,
//		o.ObjectID, o.PlotID, o.Name, o.ObjectType, o.Geom,
//	).Scan(&o.CreatedAt, &o.UpdatedAt)
//}
//
//func (r *Repository) ListObjectsByPlot(ctx context.Context, plotID string) ([]UIObject, error) {
//	rows, err := r.db.Query(ctx,
//		`SELECT object_id, plot_id, name, object_type, ST_AsGeoJSON(geom), created_at, updated_at
//		 FROM objects_ui WHERE plot_id=$1 AND is_deleted=FALSE
//		 ORDER BY created_at`,
//		plotID)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var objects []UIObject
//	for rows.Next() {
//		var o UIObject
//		var geom string
//		if err := rows.Scan(&o.ObjectID, &o.PlotID, &o.Name, &o.ObjectType, &geom,
//			&o.CreatedAt, &o.UpdatedAt); err != nil {
//			return nil, err
//		}
//		o.Geom = geom
//		objects = append(objects, o)
//	}
//	return objects, rows.Err()
//}
//
//func (r *Repository) SoftDeleteObject(ctx context.Context, objectID, plotID string) error {
//	tag, err := r.db.Exec(ctx,
//		`UPDATE objects_ui SET is_deleted=TRUE, updated_at=NOW()
//		 WHERE object_id=$1 AND plot_id=$2 AND is_deleted=FALSE`,
//		objectID, plotID)
//	if err != nil {
//		return err
//	}
//	if tag.RowsAffected() == 0 {
//		return apperrors.ErrNotFound
//	}
//	return nil
//}
//
//// =====================================================================
//// GRID & CELLS
//// =====================================================================
//
//func (r *Repository) InitializeGrid(ctx context.Context, plotID string, cols, rows int, cellSize float64, originX, originY float64) error {
//	// Генерируем все ячейки сетки
//	query := `
//		INSERT INTO grid_cells (plot_id, x_index, y_index, geom)
//		SELECT $1, x, y,
//		       ST_MakeEnvelope(
//		           $2 + (x * $4),
//		           $3 + (y * $4),
//		           $2 + ((x + 1) * $4),
//		           $3 + ((y + 1) * $4),
//		           4326
//		       )
//		FROM generate_series(0, $5 - 1) AS x
//		CROSS JOIN generate_series(0, $6 - 1) AS y
//	`
//	_, err := r.db.Exec(ctx, query, plotID, originX, originY, cellSize, cols, rows)
//	return err
//}
//
//func (r *Repository) GetGridCellsByPlot(ctx context.Context, plotID string) ([]GridCell, error) {
//	rows, err := r.db.Query(ctx,
//		`SELECT plot_id, x_index, y_index, current_bed_id, current_crop_id, is_occupied, is_shaded
//		 FROM grid_cells WHERE plot_id=$1 AND is_deleted=FALSE
//		 ORDER BY x_index, y_index`,
//		plotID)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var cells []GridCell
//	for rows.Next() {
//		var c GridCell
//		if err := rows.Scan(&c.PlotID, &c.XIndex, &c.YIndex, &c.CurrentBedID,
//			&c.CurrentCropID, &c.IsOccupied, &c.IsShaded); err != nil {
//			return nil, err
//		}
//		cells = append(cells, c)
//	}
//	return cells, rows.Err()
//}
//
//func (r *Repository) GetCellIDsInPolygon(ctx context.Context, plotID string, geomJSON string) ([]string, error) {
//	rows, err := r.db.Query(ctx,
//		`SELECT x_index, y_index FROM grid_cells
//		 WHERE plot_id=$1 AND is_deleted=FALSE
//		   AND ST_Intersects(geom, ST_GeomFromGeoJSON($2))`,
//		plotID, geomJSON)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var cellIDs []string
//	for rows.Next() {
//		var x, y int
//		if err := rows.Scan(&x, &y); err != nil {
//			return nil, err
//		}
//		cellIDs = append(cellIDs, fmt.Sprintf("%d_%d", x, y))
//	}
//	return cellIDs, rows.Err()
//}
//
//func (r *Repository) UpdateCellsForBed(ctx context.Context, plotID string, cellIDs []string, bedID string) error {
//	if len(cellIDs) == 0 {
//		return nil
//	}
//
//	// Парсим cellIDs в массив кортежей (x, y)
//	coords := make([][]int, len(cellIDs))
//	for i, cid := range cellIDs {
//		var x, y int
//		_, err := fmt.Sscanf(cid, "%d_%d", &x, &y)
//		if err != nil {
//			return fmt.Errorf("invalid cell_id format: %s", cid)
//		}
//		coords[i] = []int{x, y}
//	}
//
//	// Batch update
//	batch := &pgx.Batch{}
//	for _, coord := range coords {
//		batch.Queue(
//			`UPDATE grid_cells SET current_bed_id=$1, is_occupied=TRUE, updated_at=NOW()
//			 WHERE plot_id=$2 AND x_index=$3 AND y_index=$4`,
//			bedID, plotID, coord[0], coord[1])
//	}
//
//	results := r.db.SendBatch(ctx, batch)
//	defer results.Close()
//
//	for i := 0; i < len(coords); i++ {
//		if _, err := results.Exec(); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func (r *Repository) ClearCellsForBed(ctx context.Context, plotID, bedID string) error {
//	_, err := r.db.Exec(ctx,
//		`UPDATE grid_cells SET current_bed_id=NULL, is_occupied=FALSE, updated_at=NOW()
//		 WHERE plot_id=$1 AND current_bed_id=$2`,
//		plotID, bedID)
//	return err
//}
//
//// =====================================================================
//// PLANTINGS
//// =====================================================================
//
//func (r *Repository) PlantCrop(ctx context.Context, plotID, bedID string, cropID int, plantDate time.Time, cellIDs []string) error {
//	tx, err := r.db.Begin(ctx)
//	if err != nil {
//		return err
//	}
//	defer tx.Rollback(ctx)
//
//	// 1. Обновляем grid_cells
//	for _, cid := range cellIDs {
//		var x, y int
//		fmt.Sscanf(cid, "%d_%d", &x, &y)
//		_, err := tx.Exec(ctx,
//			`UPDATE grid_cells SET current_crop_id=$1, is_occupied=TRUE, updated_at=NOW()
//			 WHERE plot_id=$2 AND x_index=$3 AND y_index=$4`,
//			cropID, plotID, x, y)
//		if err != nil {
//			return err
//		}
//	}
//
//	// 2. Создаём записи в cell_crop_history
//	for _, cid := range cellIDs {
//		var x, y int
//		fmt.Sscanf(cid, "%d_%d", &x, &y)
//		_, err := tx.Exec(ctx,
//			`INSERT INTO cell_crop_history (history_id, plot_id, x_index, y_index, crop_id, family_id, bed_id, plant_date)
//			 SELECT $1, $2, $3, $4, $5, cf.family_id, $6, $7
//			 FROM crops c JOIN crop_families cf ON c.family_id = cf.id
//			 WHERE c.id=$5`,
//			uuid.New().String(), plotID, x, y, cropID, bedID, plantDate)
//		if err != nil {
//			return err
//		}
//	}
//
//	return tx.Commit(ctx)
//}
//
//func (r *Repository) HarvestCrop(ctx context.Context, plotID, bedID string, harvestDate time.Time, yieldKg *float64, cellIDs []string) error {
//	tx, err := r.db.Begin(ctx)
//	if err != nil {
//		return err
//	}
//	defer tx.Rollback(ctx)
//
//	// 1. Обновляем grid_cells
//	for _, cid := range cellIDs {
//		var x, y int
//		fmt.Sscanf(cid, "%d_%d", &x, &y)
//		_, err := tx.Exec(ctx,
//			`UPDATE grid_cells SET current_crop_id=NULL, is_occupied=FALSE, updated_at=NOW()
//			 WHERE plot_id=$1 AND x_index=$2 AND y_index=$3`,
//			plotID, x, y)
//		if err != nil {
//			return err
//		}
//	}
//
//	// 2. Обновляем cell_crop_history (закрываем открытые записи)
//	for _, cid := range cellIDs {
//		var x, y int
//		fmt.Sscanf(cid, "%d_%d", &x, &y)
//		_, err := tx.Exec(ctx,
//			`UPDATE cell_crop_history
//			 SET harvest_date=$1, yield_kg=$2
//			 WHERE plot_id=$3 AND x_index=$4 AND y_index=$5
//			   AND harvest_date IS NULL`,
//			harvestDate, yieldKg, plotID, x, y)
//		if err != nil {
//			return err
//		}
//	}
//
//	return tx.Commit(ctx)
//}
//
//// =====================================================================
//// EVENT STORE
//// =====================================================================
//
//func (r *Repository) AppendEvent(ctx context.Context, plotID string, eventType int, payload interface{}) error {
//	payloadJSON, err := json.Marshal(payload)
//	if err != nil {
//		return err
//	}
//
//	// Получаем следующий sequence_number
//	var nextSeq int64
//	err = r.db.QueryRow(ctx,
//		`SELECT COALESCE(MAX(sequence_number), 0) + 1 FROM event_store WHERE plot_id=$1`,
//		plotID).Scan(&nextSeq)
//	if err != nil {
//		return err
//	}
//
//	_, err = r.db.Exec(ctx,
//		`INSERT INTO event_store (event_id, plot_id, event_type, payload, sequence_number)
//		 VALUES ($1, $2, $3, $4, $5)`,
//		uuid.New().String(), plotID, eventType, payloadJSON, nextSeq)
//	return err
//}
//
//func (r *Repository) GetTimeline(ctx context.Context, plotID string, filter TimelineFilter) ([]TimelineEvent, error) {
//	query := `SELECT event_id, event_type, occurred_at,
//	                 payload->>'display_title' as display_title,
//	                 COALESCE(payload->>'display_description', '') as display_description,
//	                 COALESCE(payload->'affected_cells', '[]'::jsonb) as affected_cells
//	          FROM event_store
//	          WHERE plot_id=$1`
//	args := []interface{}{plotID}
//	argIdx := 2
//
//	if filter.From != nil {
//		query += fmt.Sprintf(" AND occurred_at >= $%d", argIdx)
//		args = append(args, *filter.From)
//		argIdx++
//	}
//	if filter.To != nil {
//		query += fmt.Sprintf(" AND occurred_at <= $%d", argIdx)
//		args = append(args, *filter.To)
//		argIdx++
//	}
//
//	query += " ORDER BY occurred_at DESC"
//
//	limit := 100
//	if filter.Limit > 0 && filter.Limit <= 500 {
//		limit = filter.Limit
//	}
//	query += fmt.Sprintf(" LIMIT $%d", argIdx)
//	args = append(args, limit)
//
//	rows, err := r.db.Query(ctx, query, args...)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var events []TimelineEvent
//	for rows.Next() {
//		var e TimelineEvent
//		var affectedCellsJSON []byte
//		if err := rows.Scan(&e.EventID, &e.EventType, &e.OccurredAt,
//			&e.DisplayTitle, &e.DisplayDesc, &affectedCellsJSON); err != nil {
//			return nil, err
//		}
//		_ = json.Unmarshal(affectedCellsJSON, &e.AffectedCells)
//		events = append(events, e)
//	}
//	return events, rows.Err()
//}
