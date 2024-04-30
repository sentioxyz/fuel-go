package main

import (
	"github.com/graph-gophers/graphql-go/types"
	"github.com/sentioxyz/fuel-go/util"
)

type generatorExtendJSON struct {
}

func (g generatorExtendJSON) Imports() []string {
	return []string{
		"encoding/json",
	}
}

func (g generatorExtendJSON) GenScalar(out *util.Output, scalarType *types.ScalarTypeDefinition, scalarGoType string) {
	switch scalarGoType {
	case "uint8", "uint16", "uint32", "uint64":
		out.Outf(`
func (s *#{name}) UnmarshalJSON(raw []byte) error {
	if i, err := unmarshalJSONUInt(raw); err != nil {
		return err
	} else {
		*s = #{name}(i)
		return nil
	}
}

func (s #{name}) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
`, "name", scalarType.TypeName())
	case "int8", "int16", "int32", "int64":
		out.Outf(`
func (s *#{name}) UnmarshalJSON(raw []byte) error {
	if i, err := unmarshalJSONInt(raw); err != nil {
		return err
	} else {
		*s = #{name}(i)
		return nil
	}
}

func (s #{name}) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
`, "name", scalarType.TypeName())
	case "float32", "float64":
		out.Outf(`
func (s *#{name}) UnmarshalJSON(raw []byte) error {
	if f, err := unmarshalJSONFloat(raw); err != nil {
		return err
	} else {
		*s = #{name}(f)
		return nil
	}
}

func (s #{name}) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
`, "name", scalarType.TypeName())
	}
}

func (g generatorExtendJSON) GenEnum(out *util.Output, enumType *types.EnumTypeDefinition) {
	out.Outf(`
func (e *#{name}) UnmarshalJSON(raw []byte) error {
	var val string
	if err := json.Unmarshal(raw, &val); err != nil {
		return err
	}
	for _, v := range #{name}Values {
		if v == val {
			*e = #{name}(val)
			return nil
		}
	}
	return fmt.Errorf("invalid value %q for enum type #{name}", val)
}
`, "name", enumType.TypeName())
}

func (g generatorExtendJSON) GenObject(out *util.Output, objectType *types.ObjectTypeDefinition) {}

func (g generatorExtendJSON) GenUnion(out *util.Output, unionType *types.Union) {
	out.Outf(`
func (u *#{name}) UnmarshalJSON(raw []byte) error {
	return unmarshalJSONUnion(raw, u)
}

func (u #{name}) MarshalJSON() ([]byte, error) {
	return marshalJSONUnion(u)
}
`, "name", unionType.TypeName())
}
