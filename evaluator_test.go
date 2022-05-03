package monkey_interpreter

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10}}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}
func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

// helper functions ----------------------------------------------------------

func testBooleanObject(t *testing.T, obj Object, expected bool) bool {
	result, ok := obj.(*BooleanObject)
	if !ok {
		t.Errorf("object is not BooleanObject. got=%T (%+v)", obj, obj)
		return false
	}
	require.Equal(t, result.Value, expected)
	return true
}

func testEval(input string) Object {
	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	return Eval(program)
}

func testIntegerObject(t *testing.T, obj Object, expected int64) bool {
	result, ok := obj.(*Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	require.Equal(t, result.Value, expected)
	return true
}
