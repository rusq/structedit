// Package structedit allows to interactively edit structure values.
//
// It it basically a wrapper around survey library that parses the structure
// tags and generates necessary "questions".
package structedit

import (
	"math"
	"reflect"
	"testing"
)

func ptr[T any](x T) *T {
	return &x
}

// supTypes - supported types.
type supTypes struct {
	B    bool
	C64  complex64
	C128 complex128
	F32  float32
	F64  float64
	I    int
	I8   int8
	I16  int16
	I32  int32
	I64  int64
	U    uint
	U8   uint8
	U16  uint16
	U32  uint32
	U64  uint64
	S    string
}

func Test_setValue(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		srcStruct  supTypes
		args       args
		wantStruct supTypes
		wantErr    bool
	}{
		{
			"bool",
			supTypes{B: false},
			args{"true"},
			supTypes{B: true},
			false,
		},
		{
			"complex64",
			supTypes{C64: complex64(2 + 4i)},
			args{"5+2i"},
			supTypes{C64: complex64(5 + 2i)},
			false,
		},
		{
			"complex128",
			supTypes{C128: complex128(2 + 4i)},
			args{"5+2i"},
			supTypes{C128: complex128(5 + 2i)},
			false,
		},
		{
			"float32",
			supTypes{F32: 1609.1991},
			args{"123.4567890123"},
			supTypes{F32: 123.4567890123},
			false,
		},
		{
			"float64",
			supTypes{F64: 1609.1991},
			args{"123.4567890123"},
			supTypes{F64: 123.4567890123},
			false,
		},
		{
			"int",
			supTypes{I: 65536},
			args{"32768"},
			supTypes{I: 32768},
			false,
		},
		{
			"int8",
			supTypes{I8: 64},
			args{"127"},
			supTypes{I8: 127},
			false,
		},
		{
			"int16",
			supTypes{I16: 32767},
			args{"-32768"},
			supTypes{I16: -32768},
			false,
		},
		{
			"int32",
			supTypes{I32: -1},
			args{"2147483647"},
			supTypes{I32: math.MaxInt32},
			false,
		},
		{
			"int64",
			supTypes{I64: -1},
			args{"9223372036854775807"},
			supTypes{I64: math.MaxInt64},
			false,
		},
		{
			"uint",
			supTypes{U: 65536},
			args{"32768"},
			supTypes{U: 32768},
			false,
		},
		{
			"uint8",
			supTypes{U8: 8},
			args{"255"},
			supTypes{U8: 255},
			false,
		},
		{
			"int16",
			supTypes{U16: 16},
			args{"65535"},
			supTypes{U16: 65535},
			false,
		},
		{
			"int32",
			supTypes{U32: 32},
			args{"4294967295"},
			supTypes{U32: math.MaxUint32},
			false,
		},
		{
			"int64",
			supTypes{U64: 64},
			args{"18446744073709551615"},
			supTypes{U64: math.MaxUint64},
			false,
		},
		{
			"string",
			supTypes{S: "before"},
			args{"after"},
			supTypes{S: "after"},
			false,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := reflect.ValueOf(&tt.srcStruct).Elem().Field(i)
			if err := setValue(val, tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("setValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.srcStruct, tt.wantStruct) {
				t.Errorf("mismatch:\nwant=%#v\ngot= %#v", tt.wantStruct, tt.srcStruct)
			}
		})
	}
}
