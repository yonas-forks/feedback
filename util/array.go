package util

func ToAnyArray(rows []map[string]any) []any {
	list := []any{}
	for _, row := range rows {
		list = append(list, row)
	}
	return list
}
