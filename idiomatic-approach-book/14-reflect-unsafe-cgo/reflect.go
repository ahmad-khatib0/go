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

func valueOf() {
	s := []string{"a", "b", "c"}
	sv := reflect.ValueOf(s)        // sv is of type reflect.Value
	s2 := sv.Interface().([]string) // s2 is of type []string
	fmt.Println(s2)

	i := 10
	iv := reflect.ValueOf(&i)
	ivv := iv.Elem()
	ivv.SetInt(20)
	fmt.Println(i) // prints 20
}

func changeIntReflect(i *int) {
	iv := reflect.ValueOf(i)
	iv.Elem().SetInt(20)
}

func makingNewValues() {
	// a trick that lets you create a variable to represent a reflect.Type if you don’t have a value handy
	// stringType       contains a reflect.Type that represents a string
	// stringSliceType  contains a reflect.Type that represents a []string
	var stringType = reflect.TypeOf((*string)(nil)).Elem() // converting nil to a pointer to string
	// without the parenthesis around string, the compiler thinks that we are converting nil to string, which is illegal.
	var stringSliceType = reflect.TypeOf([]string(nil))
	// since nil is a valid value for a slice. All we have to do is
	// type conversion of nil to a []string and pass that to reflect.Type.

	ssv := reflect.MakeSlice(stringSliceType, 0, 10)
	sv := reflect.New(stringType).Elem()
	sv.SetString("hello")
	ssv = reflect.Append(ssv, sv)
	ss := ssv.Interface().([]string)
	fmt.Println(ss) // prints [hello]
}
