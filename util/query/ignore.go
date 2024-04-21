package query

import (
	"reflect"
)

type IgnoreChecker func(object reflect.Type, field reflect.StructField) bool

func IgnoreObjects(objs ...any) IgnoreChecker {
	return func(_ reflect.Type, field reflect.StructField) bool {
		fieldType := field.Type
		for {
			for _, obj := range objs {
				if reflect.TypeOf(obj) == fieldType {
					return true
				}
			}
			switch fieldType.Kind() {
			case reflect.Pointer, reflect.Slice, reflect.Array:
				fieldType = fieldType.Elem()
			default:
				return false
			}
		}
	}
}

func IgnoreField(obj any, fieldName string) IgnoreChecker {
	return func(object reflect.Type, field reflect.StructField) bool {
		return object == reflect.TypeOf(obj) && field.Name == fieldName
	}
}

func IgnoreOtherField(obj any, fieldNames ...string) IgnoreChecker {
	return func(object reflect.Type, field reflect.StructField) bool {
		if object != reflect.TypeOf(obj) {
			return false
		}
		for _, fieldName := range fieldNames {
			if field.Name == fieldName {
				return false
			}
		}
		return true
	}
}

func MergeIgnores(checkers ...IgnoreChecker) IgnoreChecker {
	return func(object reflect.Type, field reflect.StructField) bool {
		for _, checker := range checkers {
			if checker(object, field) {
				return true
			}
		}
		return false
	}
}
