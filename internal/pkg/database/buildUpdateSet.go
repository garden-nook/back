package database

import (
	"fmt"
	"strings"
)

// SetField описывает одно обновляемое поле.
type SetField struct {
	Name  string // Имя колонки в БД
	Value any    // Значение (сам указатель, например *string)
}

// BuildUpdateSet собирает SQL‑выражение для SET‑части UPDATE и слайс значений.
// Пропускает поля, где Value == nil или указатель на nil.
func BuildUpdateSet(startIdx int, fields ...SetField) (setClause string, args []any) {
	var clauses []string
	idx := startIdx
	for _, f := range fields {
		if f.Value == nil {
			continue
		}
		clauses = append(clauses, fmt.Sprintf("%s=$%d", f.Name, idx))
		args = append(args, f.Value)
		idx++
	}
	if len(clauses) == 0 {
		return "", nil
	}
	return strings.Join(clauses, ", "), args
}
