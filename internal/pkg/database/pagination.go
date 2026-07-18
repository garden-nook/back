package database

import "fmt"

type Pagination struct {
	Limit  int
	Offset int
}

func (p Pagination) SQL(startIdx int) (string, []any) {
	if p.Limit <= 0 {
		return "", nil
	}
	offset := p.Offset
	if offset < 0 {
		offset = 0
	}
	return fmt.Sprintf(" LIMIT $%d OFFSET $%d", startIdx, startIdx+1), []any{p.Limit, offset}
}
