package models

import "strings"

func buildInClause[T comparable](vals []T) (string, []any) {
	placeholders := strings.Repeat("?,", len(vals)-1) + "?"

	args := make([]any, len(vals))
	for i, v := range vals {
		args[i] = v
	}

	return placeholders, args
}
