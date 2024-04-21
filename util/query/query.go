package query

import (
	"bytes"
	"github.com/sentioxyz/fuel-go/util"
	"reflect"
)

func (b Builder) indent() Builder {
	return Builder{
		Prefix: b.Prefix + b.Indent,
		Indent: b.Indent,
		EOL:    b.EOL,
	}
}

func (b Builder) genObjectQuery(objType reflect.Type, ignoreChecker IgnoreChecker) string {
	var buf bytes.Buffer
	var w = util.Output{Writer: &buf}

	if _, has := objType.FieldByName(util.UnionTypeFieldName); has {
		// is an union
		w.Out("%s__typename%s", b.Prefix, b.EOL)
		for i := 0; i < objType.NumField(); i++ {
			field := objType.Field(i)
			if field.Name == util.UnionTypeFieldName {
				continue
			}
			if ignoreChecker != nil && ignoreChecker(objType, field) {
				continue
			}
			w.Out("%s... on %s {%s", b.Prefix, field.Name, b.EOL)
			w.Out(b.indent().genObjectQuery(util.UnwrapGoType(field.Type), ignoreChecker))
			w.Out("%s}%s", b.Prefix, b.EOL)
		}
		return buf.String()
	}

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		name, has := field.Tag.Lookup("json")
		if !has {
			continue
		}
		if ignoreChecker != nil && ignoreChecker(objType, field) {
			continue
		}
		switch field.Tag.Get("kind") {
		case "OBJECT", "UNION":
			w.Out("%s%s {%s", b.Prefix, name, b.EOL)
			w.Out(b.indent().genObjectQuery(util.UnwrapGoType(field.Type), ignoreChecker))
			w.Out("%s}%s", b.Prefix, b.EOL)
		default:
			w.Out("%s%s%s", b.Prefix, name, b.EOL)
		}
	}
	return buf.String()
}

func (b Builder) GenObjectQuery(obj any, ignoreChecker IgnoreChecker) string {
	return b.genObjectQuery(util.UnwrapGoType(reflect.TypeOf(obj)), ignoreChecker)
}
