package structedit

import (
	"reflect"
	"testing"
)

func Test_parseTag(t *testing.T) {
	type args struct {
		annotation string
		sep        string
	}
	tests := []struct {
		name string
		args args
		want tagInfo
	}{
		{"name and descr", args{"Foo:that's a foo", ":"}, tagInfo{name: "Foo", description: "that's a foo"}},
		{"name and descr, omitempty", args{"Foo,omitempty:that's a foo", ":"}, tagInfo{name: "Foo", description: "that's a foo", omitempty: true}},
		{"skip", args{"-", ":"}, tagInfo{name: "-", skip: true}},
		{"empty", args{"", ":"}, tagInfo{name: ""}},
		{"omitempty, no name", args{",omitempty", ":"}, tagInfo{omitempty: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseTag(tt.args.annotation, tt.args.sep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isExported(t *testing.T) {
	type args struct {
		fieldName string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"unexported", args{"test"}, false},
		{"exported", args{"Test"}, true},
		{"exported UTF", args{"Нет"}, true},
		{"unexported UTF", args{"войне"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isExported(tt.args.fieldName); got != tt.want {
				t.Errorf("isExported() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testStruct struct {
	StrFoo     string `ed:"Foo:that's a foo"`
	IntBar     int    `ed:"Bar:enter a bar value"`
	Bool       bool   `ed:"bool:enter a boolean value"`
	SkipMe     string `ed:"-:this should be skipped"`
	OmitMe     int64  `ed:",omitempty:this should be skipped, when zero"`
	FloatNoTag float32
	unexported string `ed:"unexported:this should not be visible"`
}

var testVar testStruct

var testVarFilled = testStruct{
	StrFoo:     "yay",
	IntBar:     42,
	Bool:       true,
	SkipMe:     "fuck off, field, we hate you",
	OmitMe:     4,
	FloatNoTag: 2.3,
	unexported: "you can't see me",
}

func TestEditor_parseStruct(t *testing.T) {
	type fields struct {
		tagName string
		sep     string
		menu    sysmenu
	}
	type args struct {
		a any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]editField
	}{
		{
			"sample structure",
			fields{
				tagName: "ed",
				sep:     ":",
			},
			args{&testVar},
			map[string]editField{
				"StrFoo": {
					Name:        "Foo",
					Description: "that's a foo",
					Index:       0,
					Value:       reflect.ValueOf(&testVar).Elem().Field(0),
				},
				"IntBar": {
					Name:        "Bar",
					Description: "enter a bar value",
					Index:       1,
					Value:       reflect.ValueOf(&testVar).Elem().Field(1),
				},
				"Bool": {
					Name:        "bool",
					Description: "enter a boolean value",
					Index:       2,
					Value:       reflect.ValueOf(&testVar).Elem().Field(2),
				},
				"FloatNoTag": {
					Name:        "FloatNoTag",
					Description: "",
					Index:       5,
					Value:       reflect.ValueOf(&testVar).Elem().Field(5),
				},
			},
		},
		{
			"ensure that omitempty fields appear",
			fields{
				tagName: "ed",
				sep:     ":",
			},
			args{&testVarFilled},
			map[string]editField{
				"StrFoo": {
					Name:        "Foo",
					Description: "that's a foo",
					Index:       0,
					Value:       reflect.ValueOf(&testVarFilled).Elem().Field(0),
				},
				"IntBar": {
					Name:        "Bar",
					Description: "enter a bar value",
					Index:       1,
					Value:       reflect.ValueOf(&testVarFilled).Elem().Field(1),
				},
				"Bool": {
					Name:        "bool",
					Description: "enter a boolean value",
					Index:       2,
					Value:       reflect.ValueOf(&testVarFilled).Elem().Field(2),
				},
				"OmitMe": {
					Name:        "OmitMe",
					Description: "this should be skipped, when zero",
					Index:       4,
					Value:       reflect.ValueOf(&testVarFilled).Elem().Field(4),
				},
				"FloatNoTag": {
					Name:        "FloatNoTag",
					Description: "",
					Index:       5,
					Value:       reflect.ValueOf(&testVarFilled).Elem().Field(5),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := &Editor{
				tagName: tt.fields.tagName,
				sep:     tt.fields.sep,
				menu:    tt.fields.menu,
			}
			if got := ed.parseStruct(tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Editor.parseStruct() = \n%#v, want \n%#v", got, tt.want)
			}
		})
	}
}
