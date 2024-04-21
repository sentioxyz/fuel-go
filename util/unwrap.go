package util

import "reflect"

func UnwrapGoType(typ reflect.Type) reflect.Type {
	for {
		switch typ.Kind() {
		case reflect.Pointer, reflect.Slice:
			typ = typ.Elem()
		default:
			return typ
		}
	}
}

func GetPointer[V any](raw V) *V {
	return &raw
}
