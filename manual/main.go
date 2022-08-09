package main

import (
	"fmt"
	"log"

	"github.com/rusq/structedit"
)

func main() {
	var x testStruct

	print(x)

	ed := structedit.New()
	if err := ed.Ask("sample", &x); err != nil {
		log.Fatal(err)
	}
	print(x)
}

func print(a any) {
	fmt.Printf("%#v\n", a)
}

type testStruct struct {
	StrFoo     string `ed:"Foo:that's a foo"`
	IntBar     int    `ed:"Bar:enter a bar value"`
	Bool       bool   `ed:"bool:enter a boolean value"`
	SkipMe     string `ed:"-:this should be skipped"`
	OmitMe     int64  `ed:",omitempty:this should be skipped when zero"`
	FloatNoTag float32
	unexported string `ed:"unexported:this should not be visible"`
}
