package structedit

import (
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/rusq/dlog"
)

func (ed *Editor) parseStruct(a any) map[string]editField {
	var choices = map[string]editField{}

	ps := reflect.ValueOf(a)
	if ps.Kind() != reflect.Ptr {
		dlog.Panicf("must be a pointer: %v", reflect.TypeOf(a))
	}
	v := ps.Elem()
	if v.Kind() != reflect.Struct {
		dlog.Panicf("must be struct: %s", reflect.TypeOf(v))
	}

	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := typ.Field(i) // struct field
		value := v.Field(i)   // field value

		// skipping pointers and interfaces
		if field.Type.Kind() == reflect.Pointer || field.Type.Kind() == reflect.Interface {
			continue
		}

		dlog.Debugf("%q: %v\n", field.Name, value.CanSet())

		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			// flatten anonymous embeds
			nested := ed.parseStruct(value.Interface())
			for k, v := range nested {
				choices[k] = v
			}
			continue
		}
		if !isExported(field.Name) {
			continue
		}

		ti := parseTag(field.Tag.Get(ed.tagName), ed.sep)
		if ti.skip {
			continue
		}
		if ti.name == "" {
			ti.name = field.Name
		}
		if ti.omitempty && value.IsZero() {
			continue
		}

		choices[field.Name] = editField{
			Name:        ti.name,
			Index:       i,
			Description: ti.description,
			Value:       value,
		}
	}
	return choices
}

type tagInfo struct {
	name        string
	description string
	skip        bool
	omitempty   bool
}

// parseTag parses a structure tag.
func parseTag(annotation string, sep string) tagInfo {
	head, descr, _ := strings.Cut(annotation, sep)
	name, param, _ := strings.Cut(head, ",")

	return tagInfo{
		name:        name,
		description: descr,
		skip:        name == "-",
		omitempty:   param == "omitempty",
	}
}

// isExported returns true if the struct field is exported.
func isExported(fieldName string) bool {
	firstRune, _ := utf8.DecodeRuneInString(fieldName)
	if firstRune == utf8.RuneError {
		dlog.Panicf("is that even a field: %s ", fieldName)
	}
	if unicode.In(firstRune, unicode.Lu) {
		return true
	}
	return false
}
