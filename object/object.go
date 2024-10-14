package object

import (
	"bytes"
	"fmt"
	"go-interpreter/ast"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJECT      ObjectType = "Integer"
	BOOLEAN_OBJECT      ObjectType = "Boolean"
	ERROR_OBJECT        ObjectType = "Error"
	NULL_OBJECT         ObjectType = "Null"
	RETURN_VALUE_OBJECT ObjectType = "ReturnValue"
	FUNCTION_OBJECT     ObjectType = "Function"
	STRING_OBJECT       ObjectType = "String"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJECT
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJECT
}

type Error struct {
	Message string
}

func (e *Error) Inspect() string {
	return "Error: " + e.Message
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJECT
}

type Null struct{}

func (n *Null) Inspect() string {
	return "null"
}

func (n *Null) Type() ObjectType {
	return NULL_OBJECT
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJECT
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	store := make(map[string]Object)
	return &Environment{store, nil}
}

func NewInnerEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (env *Environment) Get(name string) (Object, bool) {
	obj, ok := env.store[name]
	if !ok && env.outer != nil {
		obj, ok = env.outer.Get(name)
	}
	return obj, ok
}

func (env *Environment) Set(name string, value Object) Object {
	env.store[name] = value
	return value
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJECT
}
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := make([]string, len(f.Parameters))
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fun (")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type String struct {
	Value string
}

func (s String) Type() ObjectType {
	return STRING_OBJECT
}

func (s String) Inspect() string {
	return s.Value
}
