package query

import (
	"reflect"
)

type IgnoreChecker func(object reflect.Type, field reflect.StructField) bool

func IgnoreObjects(objs ...any) IgnoreChecker {
	objTypeSet := make(map[reflect.Type]bool)
	for _, obj := range objs {
		objTypeSet[reflect.TypeOf(obj)] = true
	}
	return func(_ reflect.Type, field reflect.StructField) bool {
		fieldType := field.Type
		for {
			if objTypeSet[fieldType] {
				return true
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
	objType := reflect.TypeOf(obj)
	return func(object reflect.Type, field reflect.StructField) bool {
		return object == objType && field.Name == fieldName
	}
}

func IgnoreOtherFields(obj any, fieldNames ...string) IgnoreChecker {
	objType := reflect.TypeOf(obj)
	fieldNameSet := make(map[string]bool)
	for _, fieldName := range fieldNames {
		fieldNameSet[fieldName] = true
	}
	return func(object reflect.Type, field reflect.StructField) bool {
		return object == objType && !fieldNameSet[field.Name]
	}
}

func MergeIgnores(checkers ...IgnoreChecker) IgnoreChecker {
	switch len(checkers) {
	case 0:
		return nil
	case 1:
		return checkers[0]
	default:
		return func(object reflect.Type, field reflect.StructField) bool {
			for _, checker := range checkers {
				if checker(object, field) {
					return true
				}
			}
			return false
		}
	}
}
