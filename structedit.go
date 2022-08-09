// Package structedit allows to interactively edit structure values.
//
// It it basically a wrapper around survey library that parses the structure
// tags and generates necessary "questions".
package structedit

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
)

type Editor struct {
	tagName string
	sep     string

	menu sysmenu
}

type sysmenu struct {
	OK menuitem
}

type menuitem struct {
	label string
	descr string
}

var defEditor = &Editor{
	tagName: "ed",
	sep:     ":",
	menu: sysmenu{
		OK: menuitem{"[ OK ]", "Finish the setup"},
	},
}

func Ask(msg string, a any) error {
	return defEditor.Ask(msg, a)
}

type Option func(*Editor)

func WithTag(tag string) Option {
	return func(e *Editor) {
		if tag != "" {
			e.tagName = tag
		}
	}
}
func WithSeparator(sep string) Option {
	return func(e *Editor) {
		if sep != "" {
			e.sep = ""
		}
	}
}

func WithOKLabel(label, description string) Option {
	return func(e *Editor) {
		if label != "" {
			e.menu.OK.label = label
		}
		e.menu.OK.label = description
	}
}

func New(opt ...Option) *Editor {
	ed := &Editor{
		tagName: defEditor.tagName,
		sep:     defEditor.sep,
		menu: sysmenu{
			OK: defEditor.menu.OK,
		},
	}
	for _, f := range opt {
		f(ed)
	}
	return ed
}

func (ed *Editor) Ask(msg string, a any) error {
	parsed := ed.parseStruct(a)

	keys := make([]string, 0, len(parsed))
	for k := range parsed {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return parsed[keys[i]].Index < parsed[keys[j]].Index
	})

	keys = append(keys, []string{ed.menu.OK.label}...)

	q := &survey.Select{
		Message: msg,
		Options: keys,
		Description: func(value string, index int) string {
			if value == ed.menu.OK.label {
				return ed.menu.OK.descr
			}
			opt := parsed[value]
			if opt.Description != "" {
				return fmt.Sprintf("[%v]:  %s", opt.Value.Interface(), opt.Description)
			}
			return fmt.Sprintf("[%v]", opt.Value.Interface())
		},
	}
	for {
		var resp string
		if err := survey.AskOne(q, &resp); err != nil {
			return err
		}
		if resp == ed.menu.OK.label {
			break
		}
		opt := parsed[resp]
		if err := opt.Ask(); err != nil {
			return err
		}
	}
	return nil
}

type editField struct {
	Name        string
	Index       int
	Value       reflect.Value
	Description string
}

func (fld *editField) Validate(ans any) error {
	s := ans.(string)

	var err error
	switch fld.Value.Kind() {
	case reflect.Bool:
		_, err = strconv.ParseBool(s)
	case reflect.Complex64, reflect.Complex128:
		_, err = strconv.ParseComplex(s, 128)
	case reflect.Float32, reflect.Float64:
		_, err = strconv.ParseFloat(s, 64)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		_, err = strconv.ParseInt(s, 10, 64)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		_, err = strconv.ParseUint(s, 10, 64)
	}
	return err
}

// https://stackoverflow.com/questions/6395076
func (fld *editField) Ask() error {
	q := []*survey.Question{
		{
			Name: "value",
			Prompt: &survey.Input{
				Message: fmt.Sprintf("Input value for %q: ", fld.Name),
				Help:    fld.Description,
				Default: fmt.Sprintf("%v", fld.Value.Interface()),
			},
			Validate: survey.ComposeValidators(survey.Required, fld.Validate),
		},
	}
	var val string
	if err := survey.Ask(q, &val); err != nil {
		return err
	}
	if !fld.Value.IsValid() {
		return fmt.Errorf("invalid value for field: %q", fld.Name)
	}
	if !fld.Value.CanSet() {
		return fmt.Errorf("can't set the value of field %q", fld.Name)
	}
	if err := fld.Set(val); err != nil {
		return err
	}
	return nil
}

func (fld *editField) Set(s string) error {
	return setValue(fld.Value, s)
}

func setValue(val reflect.Value, s string) error {
	switch val.Kind() {
	case reflect.Bool:
		v, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		val.SetBool(v)
	case reflect.Complex64, reflect.Complex128:
		x, err := strconv.ParseComplex(s, 128)
		if err != nil {
			return err
		}
		val.SetComplex(x)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		val.SetFloat(f)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		val.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		val.SetUint(u)
	case reflect.String:
		val.SetString(s)
	default:
		return fmt.Errorf("unsupported type: %s", val.Kind())
	}
	return nil
}
