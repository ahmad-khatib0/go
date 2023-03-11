package main

func main() {
	// you can use any primitive type or compound type literal to define a concrete type
	primitiveTypes()
}

func primitiveTypes() {
	type Score int
	type Converter func(string) Score
	type TeamScores map[string]Score
}
