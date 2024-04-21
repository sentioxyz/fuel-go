package query

import (
	"bytes"
	"fmt"
	"github.com/sentioxyz/fuel-go/util"
	"reflect"
)

func (b Builder) genParam(value reflect.Value) string {
	var buf bytes.Buffer
	var w = util.Output{Writer: &buf}
	for i := 0; i < value.NumField(); i++ {
		tag := value.Type().Field(i).Tag
		name, has := tag.Lookup("name")
		if !has {
			continue
		}
		field := value.Field(i)
		if field.Kind() == reflect.Pointer && field.IsNil() {
			continue
		}
		if tag.Get("kind") == "INPUT_OBJECT" {
			w.Out("%s%s: {%s", b.Prefix, name, b.EOL)
			w.Out(b.indent().genParam(field))
			w.Out("%s}%s", b.Prefix, b.EOL)
		} else if str, is := field.Interface().(fmt.Stringer); is {
			w.Out("%s%s: %q%s", b.Prefix, name, str.String(), b.EOL)
		} else {
			w.Out("%s%s: %s%s", b.Prefix, name, field.String(), b.EOL)
		}
	}
	return buf.String()
}

func (b Builder) GenParam(param any) string {
	paramValue := reflect.ValueOf(param)
	if paramValue.Kind() == reflect.Pointer {
		paramValue = paramValue.Elem()
	}
	return b.genParam(paramValue)
}
