package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// ListQuery выполняет запрос с фильтрацией, сортировкой и опциональной пагинацией.
// T — тип строки (должен иметь теги db для pgx.RowToStructByName).
// baseQuery — основная часть запроса: "SELECT ... FROM ... JOIN ..." (без WHERE, ORDER BY, LIMIT).
// whereClause — строка с WHERE-условием (например, "WHERE p.owner_id = $1 AND p.is_deleted = FALSE").
// whereArgs — аргументы для WHERE-условия.
// orderBy — выражение для ORDER BY (без ключевых слов ORDER BY, например "p.name ASC").
// p — параметры пагинации; если nil — возвращаются все записи и total = len(result).
func ListQuery[T any](
	ctx context.Context,
	db DBTX,
	baseQuery string,
	whereClause string,
	whereArgs []any,
	orderBy string,
	p *Pagination,
) ([]T, int, error) {
	orderSQL := ""
	if orderBy != "" {
		orderSQL = " ORDER BY " + orderBy
	}

	if p == nil {
		query := baseQuery + " " + whereClause + orderSQL
		rows, err := db.Query(ctx, query, whereArgs...)
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()

		result, err := pgx.CollectRows(rows, pgx.RowToStructByName[T])
		if err != nil {
			return nil, 0, err
		}
		return result, len(result), nil
	}

	// Пагинация — batch из COUNT и данных
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM ( %s %s ) AS sub", baseQuery, whereClause)
	pagSQL, pagArgs := p.SQL(len(whereArgs) + 1)
	dataQuery := baseQuery + " " + whereClause + orderSQL + pagSQL

	batch := &pgx.Batch{}
	batch.Queue(countQuery, whereArgs...)

	dataQuery = baseQuery + " " + whereClause + orderSQL + pagSQL
	dataArgs := append(whereArgs, pagArgs...)
	batch.Queue(dataQuery, dataArgs...)

	results := db.SendBatch(ctx, batch)
	defer results.Close()

	var total int
	if err := results.QueryRow().Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := results.Query()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[T])
	if err != nil {
		return nil, 0, err
	}
	return result, total, nil
}
