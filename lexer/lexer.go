package lexer

import (
	"go-interpreter/token"
	"unicode"
)

type Lexer struct {
	input        string
	currPosition int
	readPosition int
	// todo)) change to rune for utf8 support
	currChar byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.eatWhitespace()

	switch l.currChar {
	case '=':
		peeked := l.peekChar()
		if peeked == '=' {
			tok = token.New(token.EQUALS, string(l.currChar)+string(peeked))
			l.readChar()
		} else {
			tok = token.New(token.ASSIGN, string(l.currChar))
		}
	case '+':
		tok = token.New(token.PLUS, string(l.currChar))
	case '-':
		tok = token.New(token.MINUS, string(l.currChar))
	case '!':
		peeked := l.peekChar()
		if peeked == '=' {
			tok = token.New(token.NOT_EQUALS, string(l.currChar)+string(peeked))
			l.readChar()
		} else {
			tok = token.New(token.BANG, string(l.currChar))
		}
	case '*':
		tok = token.New(token.ASTERISK, string(l.currChar))
	case '/':
		tok = token.New(token.SLASH, string(l.currChar))
	case '<':
		peeked := l.peekChar()
		if peeked == '=' {
			tok = token.New(token.LESS_EQUALS, string(l.currChar)+string(peeked))
			l.readChar()
		} else {
			tok = token.New(token.LESS_THAN, string(l.currChar))
		}
	case '>':
		peeked := l.peekChar()
		if peeked == '=' {
			tok = token.New(token.GREATER_EQUALS, string(l.currChar)+string(peeked))
			l.readChar()
		} else {
			tok = token.New(token.GREATER_THAN, string(l.currChar))
		}
	case ',':
		tok = token.New(token.COMMA, string(l.currChar))
	case ';':
		tok = token.New(token.SEMICOLON, string(l.currChar))
	case '(':
		tok = token.New(token.LEFT_PAREN, string(l.currChar))
	case ')':
		tok = token.New(token.RIGHT_PAREN, string(l.currChar))
	case '{':
		tok = token.New(token.LEFT_BRACE, string(l.currChar))
	case '}':
		tok = token.New(token.RIGHT_BRACE, string(l.currChar))
	case '[':
		tok = token.New(token.LEFT_BRACKET, string(l.currChar))
	case ']':
		tok = token.New(token.RIGHT_BRACKET, string(l.currChar))
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString('"')
	case '`':
		tok.Type = token.STRING
		tok.Literal = l.readString('`')
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		switch {
		case unicode.IsLetter(rune(l.currChar)):
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		case unicode.IsDigit(rune(l.currChar)):
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		default:
			tok = token.New(token.ILLEGAL, string(l.currChar))
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	// todo)) update for utf8 support
	if l.readPosition >= len(l.input) {
		l.currChar = 0
	} else {
		l.currChar = l.input[l.readPosition]
	}

	l.currPosition = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	// todo)) update for utf8 support
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) readIdentifier() string {
	startingPosition := l.currPosition
	for unicode.IsLetter(rune(l.currChar)) {
		l.readChar()
	}
	return l.input[startingPosition:l.currPosition]
}

func (l *Lexer) readNumber() string {
	startingPosition := l.currPosition
	for unicode.IsDigit(rune(l.currChar)) {
		l.readChar()
	}
	return l.input[startingPosition:l.currPosition]
}

func (l *Lexer) eatWhitespace() {
	for l.currChar == ' ' || l.currChar == '\t' || l.currChar == '\n' || l.currChar == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readString(deli byte) string {
	out := ""
	for {
		l.readChar()
		if l.currChar == deli || l.currChar == 0 {
			break
		}

		if l.currChar == '\\' {
			if l.peekChar() == '\n' {
				l.readChar()
				continue
			}

			l.readChar()

			switch l.currChar {
			case 'n':
				l.currChar = '\n'
			case 'r':
				l.currChar = '\r'
			case 't':
				l.currChar = '\t'
			case '"':
				l.currChar = '"'
			case '\\':
				l.currChar = '\\'
			case 0:
				break
			}
		}

		out = out + string(l.currChar)
	}
	return out
}
