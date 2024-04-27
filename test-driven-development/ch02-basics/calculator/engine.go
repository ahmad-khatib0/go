package calculator

// Engine is the mathematical logic part of the calculator.
type Engine struct{}

// Operation is the wrapper object that contains
// the operator and operand of a mathematical expression.
type Operation struct {
	Expression string
	Operator   string
	Operands   []float64
}

func (e *Engine) Add(x, y float64) float64 {
	return 0
}
