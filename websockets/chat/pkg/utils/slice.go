package utils

import "reflect"

func (u *Utils) FlattenMap(slice []any) []any {
	var result []any
	for _, v := range slice {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(v)
			for i := 0; i < s.Len(); i++ {
				result = append(result, s.Index(i).Interface())
			}
		default:
			result = append(result, v)
		}
	}

	return result
}
