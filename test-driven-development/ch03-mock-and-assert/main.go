package main

import (
	"flag"
	"log"

	"github.com/ahmad-khatib0/go/test-driven-development/ch03-mock-and-assert/calculator"
	"github.com/ahmad-khatib0/go/test-driven-development/ch03-mock-and-assert/input"
)

func main() {
	expr := flag.String("expression", "", "mathematical expression to parse")
	flag.Parse()

	engine := calculator.NewEngine()
	validator := input.NewValidator(engine.GetNumOperands(), engine.GetValidOperators())
	parser := input.NewParser(engine, validator)
	result, err := parser.ProcessExpression(*expr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(*result)
}
