package main

import (
	"fmt"
)

// make the sum function accept any parameter type by using the any interface:
// func sum2[T any](x, y T) T {
//   // implementation
// }

// The comparable interface is a pre-declared type constraint that denotes types that can be
// compared using the == operator. These types are string, bool, all numbers, and composite
// types containing comparable fields. For example, the sum function will now allow these types
// as parameters, but we will need to implement the sum logic accordingly, as the + operator will
// no longer be implemented by all these types:
// func sum[T comparable](x, y T) T {
//   // implementation
// }

// Type sets can be created using the | operator. This allows us to create constraints that contain
// multiple types without having to wrap them in a custom interface. This is the example we have
// seen in Figure 11.2, where the sum function allows the int64 and float64 types, which
// already support the + operator:
// func sum[T int64 | float64](x, y T) T {
//  // implementation
// }

// Custom type constraints can also be created as an interface and reused in our code. These are
// also declared using the | operator. For example, the NumberConstraint interface allows the int64
// and float64 types, which we can then use in the specification of the sum function, resulting in
// the same specification as the previous type sets:
type NumberConstraint interface {
	int64 | float64
}

// func sum[T NumberConstraint](x, y T) T {
//  // implementation
// }

// The ~ keyword can be used to restrict all custom types that have the same underlying type. This
// allows us to encompass custom types into our constraints. For example, the Number interface
// will now allow any int- and float64-based types:
type Number interface {
	~int64 | ~float64
}

func sum[T Number](x, y T) T {
	return x + y
}

// The constraints package defines some useful constraints that can be used together with
// your generic code. This package contains numerical and ordered constraints that you might
// find useful. For example, we can modify the sum function to accept all signed integers by
// using this package:
// func sum1[T constraints.Signed](x, y T) T {
//   // implementation
// }

func main() {
	fmt.Println(sum[int64](2, 3))
}
