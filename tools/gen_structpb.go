package main

import (
	"github.com/graph-gophers/graphql-go/types"
	"github.com/sentioxyz/fuel-go/util"
)

type generatorExtendStructpb struct {
}

func (g generatorExtendStructpb) Imports() []string {
	return []string{
		"reflect",
		"google.golang.org/protobuf/types/known/structpb",
	}
}

func (g generatorExtendStructpb) GenScalar(
	out *util.Output,
	scalarType *types.ScalarTypeDefinition,
	scalarGoType string,
) {
	switch scalarGoType {
	case "uint8", "uint16", "uint32", "uint64":
		out.Outf(`
func (s #{name}) MarshalStructpb() *structpb.Value {
	return structpb.NewStringValue(s.String())
}
`, "name", scalarType.TypeName())
	case "int8", "int16", "int32", "int64":
		out.Outf(`
func (s #{name}) MarshalStructpb() *structpb.Value {
	return structpb.NewStringValue(s.String())
}
`, "name", scalarType.TypeName())
	case "float32", "float64":
		out.Outf(`
func (s #{name}) MarshalStructpb() *structpb.Value {
	return structpb.NewStringValue(s.String())
}
`, "name", scalarType.TypeName())
	case "string":
		out.Outf(`
func (s #{name}) MarshalStructpb() *structpb.Value {
	return structpb.NewStringValue(string(s))
}
`, "name", scalarType.TypeName())
	case "bool":
		out.Outf(`
func (s #{name}) MarshalStructpb() *structpb.Value {
	return structpb.NewBoolValue((bool)(s))
}
`, "name", scalarType.TypeName())
	default:
		out.Outf(`
func (s #{name}) MarshalStructpb() *structpb.Value {
	return structpb.NewStringValue(s.String())
}
`, "name", scalarType.TypeName())
	}
}

func (g generatorExtendStructpb) GenEnum(out *util.Output, enumType *types.EnumTypeDefinition) {
	out.Outf(`
func (s #{name}) MarshalStructpb() *structpb.Value {
	return structpb.NewStringValue(string(s))
}
`, "name", enumType.TypeName())
}

func (g generatorExtendStructpb) GenObject(out *util.Output, objectType *types.ObjectTypeDefinition) {
	out.Outf(`
func (s #{name}) MarshalStructpb() *structpb.Value {
	return marshalStructpbObject(reflect.ValueOf(s), true)
}
`, "name", objectType.TypeName())
}

func (g generatorExtendStructpb) GenUnion(out *util.Output, unionType *types.Union) {
	out.Outf(`
func (s #{name}) MarshalStructpb() *structpb.Value {
	return marshalStructpbUnion(reflect.ValueOf(s))
}
`, "name", unionType.TypeName())
}
