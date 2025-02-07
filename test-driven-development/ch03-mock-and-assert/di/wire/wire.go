//go:build wireinject

package main

import w "github.com/google/wire"

var Set = w.NewSet(NewEngine, w.Bind(new(Adder), new(*Engine)), NewCalculator)

func InitCalc() *Calculator {
	w.Build(Set)
	return nil
}
