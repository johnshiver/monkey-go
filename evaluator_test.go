package monkey_interpreter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test case: %d", i), func(t *testing.T) {
			evaluated := testEval(tt.input)
			testIntegerObject(t, evaluated, tt.expected)
		})
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case: %d", i), func(t *testing.T) {
			evaluated := testEval(tt.input)
			testBooleanObject(t, evaluated, tt.expected)
		})
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
func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			evaluated := testEval(tt.input)
			integer, ok := tt.expected.(int)
			if ok {
				testIntegerObject(t, evaluated, int64(integer))
			} else {
				testNullObject(t, evaluated)
			}
		})
	}
}

func TestParseReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
         if (10 > 1) {
             if (10 > 1) {
                 return 10;
             }
             return 1;
         }`, 10},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			evaluated := testEval(tt.input)
			testIntegerObject(t, evaluated, tt.expected)
		})
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedMessage string
	}{
		{
			"int bool addition",
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"int bool addition 2",
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"negative boolean",
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"bool addition",
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"bool addition 2",
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"bool addition 3",
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"unbound identifier",
			"foobar",
			"identifier not found: foobar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			errObj, ok := evaluated.(*Error)
			if !ok {
				t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
				return
			}
			if errObj.Message != tt.expectedMessage {
				t.Errorf("wrong error message. expected=%q, got=%q",
					tt.expectedMessage, errObj.Message)
			}
		})
	}
}

func TestEvaluateLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
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
	env := NewEnvironment()
	return Eval(program, env)
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

func testNullObject(t *testing.T, obj Object) {
	require.Equal(t, NULL_OBJ, obj)
}
