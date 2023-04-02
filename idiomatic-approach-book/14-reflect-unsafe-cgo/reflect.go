package main

import (
	"fmt"
	"reflect"
)

//******************************************** Types ********************************************

func main() {

	reflectTypeOf()

	reflectElem()

	reflectOnStruct()
}

func reflectTypeOf() {
	var x int
	type Foo struct{}

	xt := reflect.TypeOf(x)
	fmt.Println(xt.Name()) // returns int

	f := Foo{}
	ft := reflect.TypeOf(f)
	fmt.Println(ft.Name()) // returns Foo

	xpt := reflect.TypeOf(&x) // . Some types, like a slice or a pointer, don’t have names
	fmt.Println(xpt.Name())   // returns an empty string
}

func reflectElem() {
	// - Some types in Go have references to other types and Elem is how to find out what the contained type is
	// - The Elem method also works for slices, maps, channels, and arrays.
	var x int
	xpt := reflect.TypeOf(&x)
	fmt.Println(xpt.Name())        // returns an empty string
	fmt.Println(xpt.Kind())        // returns reflect.Ptr
	fmt.Println(xpt.Elem().Name()) // returns "int"
	fmt.Println(xpt.Elem().Kind()) // returns reflect.Int
}

func reflectOnStruct() {

	type Foo struct {
		A int    `myTag:"value"`
		B string `myTag:"value2"`
	}

	var f Foo
	ft := reflect.TypeOf(f)
	for i := 0; i < ft.NumField(); i++ {
		curField := ft.Field(i)
		fmt.Println(curField.Name, curField.Type.Name(), curField.Tag.Get("myTag"))
	}
}

//******************************************** Values ********************************************
