package token

// todo)) switch to int or byte
type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func New(tokenType TokenType, literal string) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
	}
}

func LookupIdentifier(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENTIFIER
}

var keywords = map[string]TokenType{
	"fun":    FUNCTION,
	"var":    VAR,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENTIFIER = "IDENTIFIER"
	INT        = "INT"

	ASSIGN         = "="
	PLUS           = "+"
	MINUS          = "-"
	BANG           = "!"
	ASTERISK       = "*"
	SLASH          = "/"
	LESS_THAN      = "<"
	LESS_EQUALS    = "<="
	GREATER_THAN   = ">"
	GREATER_EQUALS = ">="

	EQUALS     = "=="
	NOT_EQUALS = "!="

	COMMA     = ","
	SEMICOLON = ";"

	LEFT_PAREN  = "("
	RIGHT_PAREN = ")"
	LEFT_BRACE  = "{"
	RIGHT_BRACE = "}"

	FUNCTION = "FUNCTION"
	VAR      = "VAR"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"

	STRING = "STRING"
)
