package main

import (
	"fmt"
	"reflect"
)

func main() {
	fmt.Println("vim-go")
	type FieldError interface {
		Tag() string

		ActualTag() string

		Namespace() string

		StructNamespace() string

		Field() string

		StructField() string

		Value() interface{}

		Param() string

		Kind() reflect.Kind

		Type() reflect.Type

		Translate(ut ut.Translator) string
	}
}
