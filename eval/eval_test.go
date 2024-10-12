package eval

import (
	"go-interpreter/lexer"
	"go-interpreter/object"
	"go-interpreter/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"420", 420},
		{"69", 69},
		{"420", 420},
		{"-420", -420},
		{"-69", -69},
		{"--420", 420},
		{"--69", 69},
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
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50}}

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

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !testBooleanObject(t, evaluated, tt.expected) {
			t.Logf("<< %q", tt.input)
		}
	}
}

func TestEvalConditionsExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", nil}, // anything expect "true" is falsy
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"!5", true}, // everything expect “true” is falsy
		{"!!5", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 420;", 420},
		{"return 420; 9;", 420},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{` if (10 > 1) { if (10 > 1) { return 10; }  return 1; } `, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !testIntegerObject(t, evaluated, tt.expected) {
			t.Logf("<< %q", tt.input)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true;", "type mismatch: Integer + Boolean"},
		{"5 + true; 5;", "type mismatch: Integer + Boolean"},
		{"-true", "invalid usage of `-` operator: -Boolean"},
		{"true + false;", "unknown operator: Boolean + Boolean"},
		{"5; true + false; 5", "unknown operator: Boolean + Boolean"},
		{"if (10 > 1) { true + false; }", "unknown operator: Boolean + Boolean"},
		{` if (10 > 1) { if (10 > 1) { return true + false; }  return 1; } `, "unknown operator: Boolean + Boolean"},
		{"foobar", "identifier not found: foobar"},
		{"if (false) { var x = 10; } x;", "identifier not found: x"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errorObject, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("object is not Error, got=%T(%v)", evaluated, evaluated)
			continue
		}
		if errorObject.Message != tt.expectedMessage {
			t.Errorf("wrong error message, got=%q, want=%q", errorObject.Message, tt.expectedMessage)
		}
	}
}

func TestVarStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var a = 5; a;", 5},
		{"var a = 5 * 5; a;", 25},
		{"var a = 5; var b = a; b;", 5},
		{"var a = 5; var b = a; var c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fun(x) { x + 2; };"

	evaluated := testEval(input)
	fun, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fun.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fun.Parameters)
	}

	if fun.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fun.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fun.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fun.Body.String())
	}
}

func TestFunctionCall(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var identity = fun(x) { x; }; identity(5);", 5},
		{"var identity = fun(x) { return x; }; identity(5);", 5},
		{"var double = fun(x) { x * 2; }; double(5);", 10},
		{"var add = fun(x, y) { x + y; }; add(5, 5);", 10},
		{"var add = fun(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fun(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	lxr := lexer.New(input)
	parsr := parser.New(lxr)
	program := parsr.ParseProgram()
	env := object.NewEnvironment()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("boolean value is wrong. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}
func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("Null object has wrong value. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
