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

	checkIfInterfaceValueIsNil()

	Filter()
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

// Use Reflection To Check If an Interface’s Value is nil
//  ▲
//  █   If you want to check if the value associated with an interface is nil,
//  █   you can do so with reflection using two methods: IsValid and IsNil
//  ▼

func hasNoValue(i interface{}) bool {
	iv := reflect.ValueOf(i)

	// The IsValid method returns true if reflect.Value holds anything other than a nil interface
	if !iv.IsValid() {
		return true
	}

	switch iv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func,
		reflect.Interface:
		return iv.IsNil()
	default:
		return false
	}
}

func checkIfInterfaceValueIsNil() {
	var a interface{}
	fmt.Println(a == nil, hasNoValue(a)) // prints true true

	var b *int
	fmt.Println(b == nil, hasNoValue(b)) // prints true true

	var c interface{} = b
	fmt.Println(c == nil, hasNoValue(c)) // prints false true

	var d int
	fmt.Println(hasNoValue(d)) // prints false

	var e interface{} = d
	fmt.Println(e == nil, hasNoValue(e)) // prints false false
}

// using reflect to make filtering function
func Filtering(slice interface{}, filter interface{}) interface{} {
	sv := reflect.ValueOf(slice)
	fv := reflect.ValueOf(filter)

	sliceLen := sv.Len()
	out := reflect.MakeSlice(sv.Type(), 0, sliceLen)
	for i := 0; i < sliceLen; i++ {
		curVal := sv.Index(i)
		values := fv.Call([]reflect.Value{curVal})
		if values[0].Bool() {
			out = reflect.Append(out, curVal)
		}
	}
	return out.Interface()
}

func Filter() {
	names := []string{"Andrew", "Bob", "Clara", "Hortense"}
	longNames := Filterering(names, func(s string) bool {
		return len(s) > 3
	}).([]string)
	fmt.Println(longNames) // [Andrew Clara Hortense]

	ages := []int{20, 50, 13}
	adults := Filterering(ages, func(age int) bool {
		return age >= 18
	}).([]int)
	fmt.Println(adults) // [20 50]
}
