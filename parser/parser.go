package parser

import (
	"fmt"
	"go-interpreter/ast"
	"go-interpreter/lexer"
	"go-interpreter/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS       // ==
	LESSGREATHER // < >
	SUM          // +
	PRODUCT      // *
	PREFIX       // -X !X
	CALL         // fun(X)
)

var precedences = map[token.TokenType]int{
	token.EQUALS:         EQUALS,
	token.NOT_EQUALS:     EQUALS,
	token.LESS_THAN:      LESSGREATHER,
	token.LESS_EQUALS:    LESSGREATHER,
	token.GREATER_THAN:   LESSGREATHER,
	token.GREATER_EQUALS: LESSGREATHER,
	token.PLUS:           SUM,
	token.MINUS:          SUM,
	token.SLASH:          PRODUCT,
	token.ASTERISK:       PRODUCT,
	token.LEFT_PAREN:     CALL,
}

type (
	prefixParseFunc func() ast.Expression
	infixParseFunc  func(ast.Expression) ast.Expression
)

type Parser struct {
	lxr    *lexer.Lexer
	errors []string

	currentToken token.Token
	peekToken    token.Token

	prefixParseFuncs map[token.TokenType]prefixParseFunc
	infixParseFuncs  map[token.TokenType]infixParseFunc
}

func New(lxr *lexer.Lexer) *Parser {
	p := &Parser{lxr: lxr, errors: []string{}}

	p.prefixParseFuncs = make(map[token.TokenType]prefixParseFunc)
	p.addPrefixFunc(token.IDENTIFIER, p.parseIdentifier)
	p.addPrefixFunc(token.INT, p.parseIntegerLiteral)
	p.addPrefixFunc(token.BANG, p.parsePrefixExpression)
	p.addPrefixFunc(token.MINUS, p.parsePrefixExpression)
	p.addPrefixFunc(token.TRUE, p.parseBoolean)
	p.addPrefixFunc(token.FALSE, p.parseBoolean)
	p.addPrefixFunc(token.LEFT_PAREN, p.parseGroupedExpression)
	p.addPrefixFunc(token.IF, p.parseIfExpression)
	p.addPrefixFunc(token.FUNCTION, p.parseFunctionLiteral)
	p.addPrefixFunc(token.STRING, p.parseStringLiteral)

	p.infixParseFuncs = make(map[token.TokenType]infixParseFunc)
	p.addInfixFunc(token.PLUS, p.parseInfixExpression)
	p.addInfixFunc(token.MINUS, p.parseInfixExpression)
	p.addInfixFunc(token.SLASH, p.parseInfixExpression)
	p.addInfixFunc(token.ASTERISK, p.parseInfixExpression)
	p.addInfixFunc(token.EQUALS, p.parseInfixExpression)
	p.addInfixFunc(token.NOT_EQUALS, p.parseInfixExpression)
	p.addInfixFunc(token.LESS_THAN, p.parseInfixExpression)
	p.addInfixFunc(token.LESS_EQUALS, p.parseInfixExpression)
	p.addInfixFunc(token.GREATER_THAN, p.parseInfixExpression)
	p.addInfixFunc(token.GREATER_EQUALS, p.parseInfixExpression)
	p.addInfixFunc(token.LEFT_PAREN, p.parseCallExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lxr.NextToken()
}

func (p *Parser) currentTokenEquals(expected token.TokenType) bool {
	return p.currentToken.Type == expected
}

func (p *Parser) peekTokenEquals(expected token.TokenType) bool {
	return p.peekToken.Type == expected
}

func (p *Parser) expectPeek(expected token.TokenType) bool {
	if p.peekTokenEquals(expected) {
		p.nextToken()
		return true
	}
	p.peekError(expected)
	return false
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %q, got %q instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFuncError(t token.TokenType) {
	p.errors = append(p.errors, fmt.Sprintf("no prefix parse function for %s found", t))
}

func (p *Parser) addPrefixFunc(tokenType token.TokenType, fun prefixParseFunc) {
	p.prefixParseFuncs[tokenType] = fun
}

func (p *Parser) addInfixFunc(tokenType token.TokenType, fun infixParseFunc) {
	p.infixParseFuncs[tokenType] = fun
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currentTokenEquals(token.EOF) {
		statement := p.parseStatement()
		program.Statements = append(program.Statements, statement)
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	statement := &ast.VarStatement{Token: p.currentToken}
	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	statement.Value = p.parseExpression(LOWEST)
	for p.peekTokenEquals(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()
	statement.ReturnValue = p.parseExpression(LOWEST)

	for p.peekTokenEquals(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: p.currentToken}
	statement.Value = p.parseExpression(LOWEST)
	if p.peekTokenEquals(token.SEMICOLON) {
		p.nextToken()
	}
	return statement
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFuncs[p.currentToken.Type]
	if prefix == nil {
		p.noPrefixParseFuncError(p.currentToken.Type)
		return nil
	}
	leftExpression := prefix()

	for !p.peekTokenEquals(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFuncs[p.peekToken.Type]
		if infix == nil {
			return leftExpression
		}

		p.nextToken()
		leftExpression = infix(leftExpression)
	}
	return leftExpression
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: p.currentToken}

	value, err := strconv.ParseInt(p.currentToken.Literal, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	literal.Value = value
	return literal
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.currentToken,
		Value: p.currentTokenEquals(token.TRUE),
	}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	expression := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RIGHT_PAREN) {
		return nil
	}
	return expression
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currentToken}
	if !p.expectPeek(token.LEFT_PAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RIGHT_PAREN) {
		return nil
	}

	if !p.expectPeek(token.LEFT_BRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenEquals(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LEFT_BRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}
	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}
	p.nextToken()

	for !p.currentTokenEquals(token.RIGHT_BRACE) && !p.currentTokenEquals(token.EOF) {
		statement := p.parseStatement()
		block.Statements = append(block.Statements, statement)
		p.nextToken()
	}
	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{Token: p.currentToken}

	if !p.expectPeek(token.LEFT_PAREN) {
		return nil
	}

	literal.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LEFT_BRACE) {
		return nil
	}

	literal.Body = p.parseBlockStatement()
	return literal
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var identifiers []*ast.Identifier

	p.nextToken()
	if p.currentTokenEquals(token.RIGHT_PAREN) {
		return identifiers
	}

	identifier := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	identifiers = append(identifiers, identifier)

	for p.peekTokenEquals(token.COMMA) {
		p.nextToken() // comma
		p.nextToken()
		identifier := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		identifiers = append(identifiers, identifier)
	}

	if !p.expectPeek(token.RIGHT_PAREN) {
		return nil
	}
	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Token: p.currentToken, Function: function}
	expression.Arguments = p.parseCallArguments()
	return expression
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	p.nextToken()
	if p.currentTokenEquals(token.RIGHT_PAREN) {
		return args
	}

	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenEquals(token.COMMA) {
		p.nextToken() // comma
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RIGHT_PAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: token.Token{},
		Value: p.currentToken.Literal,
	}
}
