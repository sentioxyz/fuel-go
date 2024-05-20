package types

import (
	"fmt"
	"google.golang.org/protobuf/types/known/structpb"
	"reflect"
)

type StructpbMarshaller interface {
	MarshalStructpb() *structpb.Value
}

func marshalStructpbObject(v reflect.Value, isObject bool) *structpb.Value {
	switch v.Kind() {
	case reflect.Pointer:
		if v.IsNil() {
			return structpb.NewNullValue()
		}
		return marshalStructpbObject(v.Elem(), isObject)
	case reflect.Slice, reflect.Array:
		listValue := structpb.ListValue{Values: make([]*structpb.Value, v.Len())}
		for i, n := 0, v.Len(); i < n; i++ {
			listValue.Values[i] = marshalStructpbObject(v.Index(i), isObject)
		}
		return structpb.NewListValue(&listValue)
	case reflect.Struct:
		if !isObject {
			if adapter, is := v.Interface().(StructpbMarshaller); is {
				return adapter.MarshalStructpb()
			}
			panic(fmt.Errorf("scalar %s is not a StructpbMarshaller", v.Type()))
		}
		structValue := structpb.Struct{Fields: make(map[string]*structpb.Value)}
		for i, n, t := 0, v.NumField(), v.Type(); i < n; i++ {
			tag := t.Field(i).Tag
			name, has := tag.Lookup("json")
			if !has {
				continue
			}
			fieldValue := marshalStructpbObject(v.Field(i), tag.Get("kind") == "OBJECT")
			if fieldValue.AsInterface() == nil {
				continue
			}
			structValue.Fields[name] = fieldValue
		}
		return structpb.NewStructValue(&structValue)
	default:
		if !isObject {
			if adapter, is := v.Interface().(StructpbMarshaller); is {
				return adapter.MarshalStructpb()
			}
			panic(fmt.Errorf("scalar %s is not a StructpbMarshaller", v.Type()))
		}
		panic(fmt.Errorf("object %s is not struct", v.Type()))
	}
}

func marshalStructpbUnion(val reflect.Value) *structpb.Value {
	vt := val.Type()
	if _, has := vt.FieldByName("TypeName_"); !has {
		panic(fmt.Errorf("%s is not an union type because miss field TypeName_", vt.Name()))
	}
	typeName := val.FieldByName("TypeName_").Interface().(string)
	if typeName == "" {
		return structpb.NewNullValue()
	}
	if _, has := vt.FieldByName(typeName); !has {
		panic(fmt.Errorf("union type %s cannot be an %q", vt.Name(), typeName))
	}
	subVal := val.FieldByName(typeName)
	if subVal.IsNil() {
		panic(fmt.Errorf("%s can not be nil", typeName))
	}
	adapter, is := subVal.Interface().(StructpbMarshaller)
	if !is {
		panic(fmt.Errorf("%s.%s is not a StructpbMarshaller", vt.Name(), typeName))
	}
	value := adapter.MarshalStructpb()
	value.GetStructValue().Fields["__typename"] = structpb.NewStringValue(typeName)
	return value
}
