package ast

import (
	"bytes"
	"go-interpreter/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type VarStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (v *VarStatement) TokenLiteral() string {
	return v.Token.Literal
}

func (v *VarStatement) String() string {
	var out bytes.Buffer

	out.WriteString(v.TokenLiteral() + " ")
	out.WriteString(v.Name.String())
	out.WriteString(" = ")

	if v.Value != nil {
		out.WriteString(v.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

func (v *VarStatement) statementNode() {
	// TODO implement me
	panic("implement me")
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

func (r *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenLiteral() + " ")
	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}

func (r *ReturnStatement) statementNode() {
	// TODO implement me
	panic("implement me")
}

type ExpressionStatement struct {
	Token token.Token
	Value Expression
}

func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}

func (e *ExpressionStatement) String() string {
	if e.Value != nil {
		return e.Value.String()
	}
	return ""
}

func (e *ExpressionStatement) statementNode() {
	// TODO implement me
	panic("implement me")
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (b *BlockStatement) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range b.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (b *BlockStatement) statementNode() {
	// TODO implement me
	panic("implement me")
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i Identifier) String() string {
	return i.Value
}

func (i Identifier) statementNode() {
	// TODO implement me
	panic("implement me")
}

func (i Identifier) expressionNode() {
	// TODO implement me
	panic("implement me")
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i IntegerLiteral) String() string {
	return i.TokenLiteral()
}

func (i IntegerLiteral) expressionNode() {
	// TODO implement me
	panic("implement me")
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (s *StringLiteral) TokenLiteral() string {
	return s.Token.Literal
}

func (s *StringLiteral) String() string {
	return s.Token.Literal
}

func (s *StringLiteral) expressionNode() {
	//TODO implement me
	panic("implement me")
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

func (p PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")

	return out.String()
}

func (p PrefixExpression) expressionNode() {
	// TODO implement me
	panic("implement me")
}

type InfixExpression struct {
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

func (i InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" ")
	out.WriteString(i.Operator)
	out.WriteString(" ")
	out.WriteString(i.Right.String())
	out.WriteString(")")

	return out.String()
}

func (i InfixExpression) expressionNode() {
	// TODO implement me
	panic("implement me")
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.TokenLiteral()
}

func (b *Boolean) expressionNode() {
	// TODO implement me
	panic("implement me")
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())
	if i.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(i.Alternative.String())
	}

	return out.String()
}

func (i *IfExpression) expressionNode() {
	// TODO implement me
	panic("implement me")
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *FunctionLiteral) TokenLiteral() string {
	return f.Token.Literal
}

func (f *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(f.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())

	return out.String()
}

func (f *FunctionLiteral) expressionNode() {
	// TODO implement me
	panic("implement me")
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (c *CallExpression) TokenLiteral() string {
	return c.Token.Literal
}

func (c *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range c.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(c.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

func (c *CallExpression) expressionNode() {
	// TODO implement me
	panic("implement me")
}
