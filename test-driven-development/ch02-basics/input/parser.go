package input

import "github.com/ahmad-khatib0/go/test-driven-development/ch02-basics/calculator"

type Parser struct {
	engine    *calculator.Engine
	validator *Validator
}
