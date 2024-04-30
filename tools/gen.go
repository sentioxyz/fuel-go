package main

import (
	"flag"
	"fmt"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/types"
	"github.com/sentioxyz/fuel-go/util"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

type Config struct {
	Package      string            `yaml:"package"`
	Import       []string          `yaml:"import"`
	ScalarMapper map[string]string `yaml:"scalarMapper"`
}

func convertToGoType(typ types.Type, wrappedNonNull bool) string {
	switch t := typ.(type) {
	case *types.NonNull:
		return convertToGoType(t.OfType, true)
	case *types.List:
		return "[]" + convertToGoType(t.OfType, false)
	default:
		if wrappedNonNull {
			return t.String()
		}
		return "*" + t.String()
	}
}

func unwrapType(typ types.Type) types.Type {
	for {
		switch t := typ.(type) {
		case *types.NonNull:
			typ = t.OfType
		case *types.List:
			typ = t.OfType
		default:
			return typ
		}
	}
}

type generator struct {
	schema      *types.Schema
	schemaTypes []types.NamedType
	extends     []GeneratorExtend
	config      Config
	w           *util.Output
}

func newGenerator(schema *types.Schema, config Config, out io.Writer, extends []GeneratorExtend) generator {
	schemaTypes := make([]types.NamedType, 0, len(schema.Types))
	for _, typ := range schema.Types {
		schemaTypes = append(schemaTypes, typ)
	}
	sort.Slice(schemaTypes, func(i, j int) bool {
		return schemaTypes[i].TypeName() < schemaTypes[j].TypeName()
	})

	return generator{
		schema:      schema,
		schemaTypes: schemaTypes,
		extends:     extends,
		config:      config,
		w:           &util.Output{Writer: out},
	}
}

type GeneratorExtend interface {
	Imports() []string
	GenScalar(out *util.Output, scalarType *types.ScalarTypeDefinition, scalarGoType string)
	GenEnum(out *util.Output, enumType *types.EnumTypeDefinition)
	GenObject(out *util.Output, objectType *types.ObjectTypeDefinition)
	GenUnion(out *util.Output, unionType *types.Union)
}

func (g generator) genHeader() {
	// package
	g.w.Out("package %s\n\n", g.config.Package)
	// import
	g.w.Out("import (\n")
	imports := map[string]bool{
		"fmt":     true,
		"strconv": true,
	}
	for _, ex := range g.extends {
		for _, imp := range ex.Imports() {
			imports[imp] = true
		}
	}
	for _, imp := range g.config.Import {
		imports[imp] = true
	}
	for imp := range imports {
		g.w.Out("\t%q\n", imp)
	}
	g.w.Out(")\n\n")
}

func (g generator) genScalars() {
	g.w.Out("// ====================\n// Scalars\n// --------------------\n\n")
	for _, typ := range g.schemaTypes {
		scalarType, is := typ.(*types.ScalarTypeDefinition)
		if !is {
			continue
		}
		scalarGoType, has := g.config.ScalarMapper[scalarType.TypeName()]
		if !has {
			log.Fatalf("miss mapping for scalar type %q", scalarType.TypeName())
		}
		switch scalarGoType {
		case "uint8", "uint16", "uint32", "uint64":
			g.w.Out("type %s %s\n", scalarType.TypeName(), scalarGoType)
			g.w.Outf(`
func (s #{name}) String() string {
	return strconv.FormatUint(uint64(s), 10)
}
`, "name", scalarType.TypeName())
		case "int8", "int16", "int32", "int64":
			g.w.Out("type %s %s\n", scalarType.TypeName(), scalarGoType)
			g.w.Outf(`
func (s #{name}) String() string {
	return strconv.FormatInt(int64(s), 10)
}
`, "name", scalarType.TypeName())
		case "float32", "float64":
			g.w.Out("type %s %s\n", scalarType.TypeName(), scalarGoType)
			g.w.Outf(`
func (s #{name}) String() string {
	return strconv.FormatFloat(float64(s), 'f', 20, 64)
}
`, "name", scalarType.TypeName())
		case "string":
			g.w.Out("type %s %s\n", scalarType.TypeName(), scalarGoType)
		case "bool":
			g.w.Out("type %s %s\n", scalarType.TypeName(), scalarGoType)
			g.w.Outf(`
func (s #{name}) String() string {
	return strconv.FormatBool(bool(s))
}
`, "name", scalarType.TypeName())
		default:
			g.w.Out("type %s struct { %s }\n", scalarType.TypeName(), scalarGoType)
		}
		for _, ex := range g.extends {
			ex.GenScalar(g.w, scalarType, scalarGoType)
		}
	}
}

func (g generator) genEnums() {
	g.w.Out("// ====================\n// Enums\n// --------------------\n\n")
	for _, typ := range g.schemaTypes {
		enumType, is := typ.(*types.EnumTypeDefinition)
		if !is || strings.HasPrefix(enumType.TypeName(), "_") {
			continue
		}
		g.w.OutLines(enumType.Desc, "// ")
		g.w.Out("type %s string\n\n", enumType.TypeName())
		g.w.Out("var %sValues = []string{\n", enumType.TypeName())
		for _, val := range enumType.EnumValuesDefinition {
			g.w.Out("  %q,\n", val.EnumValue)
		}
		g.w.Out("}\n\n")
		for _, ex := range g.extends {
			ex.GenEnum(g.w, enumType)
		}
	}
}

func (g generator) genObjects() {
	g.w.Out("// ====================\n// Objects\n// --------------------\n\n")
	for _, typ := range g.schemaTypes {
		objectType, is := typ.(*types.ObjectTypeDefinition)
		if !is || strings.HasPrefix(objectType.TypeName(), "_") {
			continue
		}
		g.w.OutLines(objectType.Desc, "// ")
		g.w.Out("type %s struct {\n", objectType.TypeName())
		for _, field := range objectType.Fields {
			goFieldName := util.UpperFirst(field.Name)
			tags := fmt.Sprintf(`json:"%s" kind:"%s"`, field.Name, unwrapType(field.Type).Kind())
			g.w.OutLines(field.Desc, "\t// ")
			g.w.Out("\t// SCHEMA: %s %s\n", field.Name, field.Type.String())
			g.w.Out("\t%s %s `%s`\n", goFieldName, convertToGoType(field.Type, false), tags)
		}
		g.w.Out("}\n")
		for _, ex := range g.extends {
			ex.GenObject(g.w, objectType)
		}
	}
}

func (g generator) genUnions() {
	g.w.Out("// ====================\n// Unions\n// --------------------\n\n")
	for _, unionType := range g.schema.Unions {
		g.w.OutLines(unionType.Desc, "// ")
		g.w.Out("type %s struct {\n", unionType.TypeName())
		g.w.Out("\t%s string `json:\"__typename\"`\n\n", util.UnionTypeFieldName)
		for _, mem := range unionType.UnionMemberTypes {
			g.w.Out("\t*%s\n", mem.TypeName())
		}
		g.w.Out("}\n\n")
		for _, ex := range g.extends {
			ex.GenUnion(g.w, unionType)
		}
	}
}

func (g generator) genInputObjects() {
	g.w.Out("// ====================\n// InputObjects\n// --------------------\n\n")
	for _, typ := range g.schemaTypes {
		inputObjectType, is := typ.(*types.InputObject)
		if !is {
			continue
		}
		g.w.OutLines(inputObjectType.Desc, "// ")
		g.w.Out("type %s struct {\n", inputObjectType.TypeName())
		for _, field := range inputObjectType.Values {
			goFieldName := util.UpperFirst(field.Name.Name)
			tags := fmt.Sprintf(`name:"%s" kind:"%s"`, field.Name.Name, unwrapType(field.Type).Kind())
			g.w.OutLines(field.Desc, "\t// ")
			g.w.Out("\t// SCHEMA: %s %s\n", field.Name.Name, field.Type.String())
			g.w.Out("\t%s %s `%s`\n", goFieldName, convertToGoType(field.Type, false), tags)
		}

		g.w.Out("}\n\n")
	}
}

func (g generator) genQueryParameters() {
	g.w.Out("// ====================\n// QueryArgumentObjects \n// --------------------\n\n")
	for {
		queryType, has := g.schema.Types["Query"]
		if !has {
			break
		}
		queryObject, is := queryType.(*types.ObjectTypeDefinition)
		if !is {
			break
		}
		for _, field := range queryObject.Fields {
			g.w.Out("type Query%sParams struct {\n", util.UpperFirst(field.Name))
			for _, arg := range field.Arguments {
				goFieldName := util.UpperFirst(arg.Name.Name)
				tags := fmt.Sprintf(`name:"%s" kind:"%s"`, arg.Name.Name, unwrapType(arg.Type).Kind())
				g.w.OutLines(arg.Desc, "\t// ")
				g.w.Out("\t// SCHEMA: %s %s\n", arg.Name.Name, arg.Type.String())
				g.w.Out("\t%s %s `%s`\n", goFieldName, convertToGoType(arg.Type, false), tags)
			}
			g.w.Out("}\n\n")
		}
		break
	}
}

func (g generator) Gen() {
	g.w.Out("// Auto generated by tools, do not edit\n\n")
	g.genHeader()
	g.genScalars()
	g.genEnums()
	g.genObjects()
	g.genUnions()
	g.genInputObjects()
	g.genQueryParameters()
}

func main() {
	schemaFile := flag.String("schema-file", "./schema.graphql", "path of the schema file")
	configFile := flag.String("config-file", "./config.yaml", "path of the config yaml file")
	outputFile := flag.String("output-file", "../types/types.go", "path of the output file")

	flag.Parse()

	// load schema
	schemaCnt, err := os.ReadFile(*schemaFile)
	if err != nil {
		log.Fatalf("read schema file failed: %v", err)
	}
	schema, err := graphql.ParseSchema(string(schemaCnt), nil)
	if err != nil {
		log.Fatalf("parse schema failed: %v", err)
	}

	// load config
	configCnt, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("read config file failed: %v", err)
	}
	var conf Config
	if err = yaml.Unmarshal(configCnt, &conf); err != nil {
		log.Fatalf("parse config failed: %v", err)
	}

	// open output file
	out, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalf("open output file failed: %v", err)
	}
	defer out.Close()

	// gen
	extends := []GeneratorExtend{
		generatorExtendJSON{},
		generatorExtendStructpb{},
	}
	newGenerator(schema.ASTSchema(), conf, out, extends).Gen()
}
