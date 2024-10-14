package lexer

import (
	"go-interpreter/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `var five = 5;
var ten = 10;

var add = fun(x, y) {
  x + y;
};

var result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;

"nice"
"hello \"world\""
`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.VAR, "var"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.VAR, "var"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fun"},
		{token.LEFT_PAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RIGHT_PAREN, ")"},
		{token.LEFT_BRACE, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RIGHT_BRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.VAR, "var"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LEFT_PAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RIGHT_PAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LESS_THAN, "<"},
		{token.INT, "10"},
		{token.GREATER_THAN, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LEFT_PAREN, "("},
		{token.INT, "5"},
		{token.LESS_THAN, "<"},
		{token.INT, "10"},
		{token.RIGHT_PAREN, ")"},
		{token.LEFT_BRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RIGHT_BRACE, "}"},
		{token.ELSE, "else"},
		{token.LEFT_BRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RIGHT_BRACE, "}"},
		{token.INT, "10"},
		{token.EQUALS, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQUALS, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.STRING, "nice"},
		{token.STRING, "hello \"world\""},
		{token.EOF, ""},
	}

	lexer := New(input)

	for i, tt := range tests {
		tok := lexer.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tok.Type wrong. expected=%q, got=%q, token=%v", i, tt.expectedType, tok.Type, tok)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - tok.Literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
