package input

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ahmad-khatib0/go/test-driven-development/ch03-mock-and-assert/calculator"
	"github.com/ahmad-khatib0/go/test-driven-development/ch03-mock-and-assert/format"
)

const expressionLength = 3

// OperationProcessor is the interface for processing mathematical operations
type OperationProcessor interface {
	ProcessOperation(operation calculator.Operation) (*string, error)
}

// ValidationHelper is the interface for input valeidation
type ValidationHelper interface {
	CheckInput(operator string, operands []float64) error
}

// Parser is responsible for converting input to mathematical operations
type Parser struct {
	engine    OperationProcessor
	validator ValidationHelper
}

// NewParser creates a ready to user parser.
func NewParser(op OperationProcessor, v ValidationHelper) *Parser {
	return &Parser{
		engine:    op,
		validator: v,
	}
}

// ProcessExpression parses an expression and sends it to the calculator
func (p *Parser) ProcessExpression(expr string) (*string, error) {
	operation, err := p.getOperation(expr)
	if err != nil {
		return nil, format.Error(expr, err)
	}
	return p.engine.ProcessOperation(*operation)
}

func (p *Parser) getOperation(expr string) (*calculator.Operation, error) {
	ops := strings.Fields(expr)
	if len(ops) != expressionLength {
		return nil, fmt.Errorf("incorrect expression length:got %d, want %d",
			len(ops), expressionLength)
	}
	leftOp, err := strconv.ParseFloat(ops[0], 64)
	if err != nil {
		return nil, fmt.Errorf("unable to process expression:%v", err)
	}
	rightOp, err := strconv.ParseFloat(ops[2], 64)
	if err != nil {
		return nil, fmt.Errorf("unable to process expression:%v", err)
	}
	operator := ops[1]
	operands := []float64{leftOp, rightOp}
	if err := p.validator.CheckInput(operator, operands); err != nil {
		return nil, err
	}

	return &calculator.Operation{
		Expression: expr,
		Operator:   operator,
		Operands:   operands,
	}, nil
}
