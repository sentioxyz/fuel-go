package types

import (
	"encoding/json"
	"fmt"
	"github.com/sentioxyz/fuel-go/util"
	"reflect"
	"strconv"
)

func unmarshalJSONUInt(raw []byte) (uint64, error) {
	if len(raw) >= 2 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
	}
	return strconv.ParseUint(string(raw), 10, 64)
}

func unmarshalJSONInt(raw []byte) (int64, error) {
	if len(raw) >= 2 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
	}
	return strconv.ParseInt(string(raw), 10, 64)
}

func unmarshalJSONFloat(raw []byte) (float64, error) {
	if len(raw) >= 2 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
	}
	return strconv.ParseFloat(string(raw), 64)
}

func unmarshalJSONUnion(raw []byte, unionObj any) error {
	var union struct {
		TypeName string `json:"__typename"`
	}
	if err := json.Unmarshal(raw, &union); err != nil {
		return err
	}
	if union.TypeName == "" {
		return nil
	}
	pv := reflect.ValueOf(unionObj)
	if pv.Kind() != reflect.Pointer || pv.IsNil() {
		return &json.InvalidUnmarshalError{Type: reflect.TypeOf(unionObj)}
	}
	rv := pv.Elem()
	rt := rv.Type()
	if _, has := rt.FieldByName(util.UnionTypeFieldName); !has {
		return fmt.Errorf("%s is not an union type because miss field #{typeFieldName}", rt.Name())
	}
	rv.FieldByName(util.UnionTypeFieldName).SetString(union.TypeName)
	for i := 0; i < rt.NumField(); i++ {
		if rt.Field(i).Name == union.TypeName {
			if rt.Field(i).Type.Kind() != reflect.Pointer {
				return fmt.Errorf("member %s of union type %T should be an pointer", union.TypeName, unionObj)
			}
			fv := reflect.New(rt.Field(i).Type.Elem())
			if err := json.Unmarshal(raw, fv.Interface()); err != nil {
				return err
			}
			rv.Field(i).Set(fv)
			return nil
		}
	}
	return fmt.Errorf("union type %T do not have member %q", unionObj, union.TypeName)
}

func marshalJSONUnion(unionObj any) ([]byte, error) {
	val := reflect.ValueOf(unionObj)
	vt := val.Type()
	if _, has := vt.FieldByName(util.UnionTypeFieldName); !has {
		return nil, fmt.Errorf("%s is not an union type because miss field %s", vt.Name(), util.UnionTypeFieldName)
	}
	typeName := val.FieldByName(util.UnionTypeFieldName).Interface().(string)
	if typeName == "" {
		return json.Marshal(nil)
	}
	if _, has := vt.FieldByName(typeName); !has {
		return nil, fmt.Errorf("union type %s cannot be an %q", vt.Name(), typeName)
	}
	subVal := val.FieldByName(typeName)
	if subVal.IsNil() {
		return nil, fmt.Errorf("%s can not be nil", typeName)
	}
	subVal = subVal.Elem()
	subTyp := subVal.Type()
	fields := make([]reflect.StructField, subVal.NumField()+1)
	fields[0], _ = vt.FieldByName(util.UnionTypeFieldName)
	for i := 0; i < subTyp.NumField(); i++ {
		fields[i+1] = subTyp.Field(i)
	}
	merged := reflect.New(reflect.StructOf(fields)).Elem()
	merged.Field(0).SetString(typeName)
	for i := 0; i < subVal.NumField(); i++ {
		merged.Field(i + 1).Set(subVal.Field(i))
	}
	return json.Marshal(merged.Interface())
}
