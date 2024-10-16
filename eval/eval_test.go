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
		{"1 <= 2", true},
		{"1 >= 2", false},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{`"nice" == "nice"`, true},
		{`"hello" != "nice"`, true},
		{`"hello" == "nice"`, false},
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
		{`"hello" - "world"`, "unknown operator: String - String"},
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

func TestStringLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"nice"`, "nice"},
		{`"hello" + " world"`, "hello world"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		str, ok := evaluated.(*object.String)
		if !ok {
			t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
		}

		if str.Value != tt.expected {
			t.Errorf("value is wrong. got=%q, want=%q", str.Value, tt.expected)
		}
	}
}

func TestBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("hello")`, 5},
		{"len(`hello`)", 5},
		{"len([1,2,3])", 3},
		{`len(69)`, object.Error{Message: "argument to `len` not supported, got Integer"}},
		{`len("one", "one")`, object.Error{Message: "wrong number of arguments, got=2, want=1"}},
		//{`puts("hello", "world!")`, nil},
		{`head([1, 2, 3])`, 1},
		{`head([])`, nil},
		{`head("hello")`, "h"},
		{`head("")`, nil},
		{`head(1)`, object.Error{Message: "argument to `head` not supported, got Integer"}},
		{`tail([1, 2, 3])`, []int{2, 3}},
		{`tail([])`, nil},
		{`tail("hello")`, "ello"},
		{`tail("")`, nil},
		{`last([1, 2, 3])`, 3},
		{`tail(1)`, object.Error{Message: "argument to `tail` not supported, got Integer"}},
		{`last([])`, nil},
		{`last("hello")`, "o"},
		{`last("")`, nil},
		{`last(1)`, object.Error{Message: "argument to `last` not supported, got Integer"}},
		{`push([], 1)`, []int{1}},
		{`push(1, 1)`, object.Error{Message: "argument to `push` must be Array, got Integer"}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case object.Error:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T(%v)", evaluated, evaluated)
				continue
			}

			if errObj.Message != expected.Message {
				t.Errorf("wrong error message, got=%q, want=%q", errObj.Message, expected)
			}
		case string:
			str, ok := evaluated.(*object.String)
			if !ok {
				t.Errorf("object is not String. got=%T(%v)", evaluated, evaluated)
				continue
			}
			if str.Value != expected {
				t.Errorf("wrong string, got=%q, want=%q", str.Value, expected)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[420, 69, 2 * 2]"

	evaluated := testEval(input)
	arr, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(arr.Elements) != 3 {
		t.Fatalf("array has wrong number of elements. got=%d, want=%d", len(arr.Elements), 3)
	}

	testIntegerObject(t, arr.Elements[0], 420)
	testIntegerObject(t, arr.Elements[1], 69)
	testIntegerObject(t, arr.Elements[2], 4)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"var i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"var myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"var myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"var myArray = [1, 2, 3]; var i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
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

func BenchmarkFib(b *testing.B) {
	input := `
var fib = fun (n) {
if (n <=1) {
return n;
}
return fib(n-1) + fib(n-2);
}

fib(20)
`
	for i := 0; i < b.N; i++ {
		testEval(input)
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
