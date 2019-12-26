package lexer

// page 18
import (
	"fmt"
	"xlang/token"
)

// Lexer basically stores the current input line and processes it.
type Lexer struct {
	input        string
	position     int  // current position
	readPosition int  // next position after current char
	ch           byte // current char
}

// New Returns a new Lexer
func New(input string) *Lexer {
	l := &Lexer{input: input}
	// Initialize to first char.
	l.readChar()
	return l
}

func (l *Lexer) readIndentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// NextToken Returns the next token of an input
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhiteSpace()
	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIndentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		}
		if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		}
		tok = newToken(token.ILLEGAL, l.ch)
	}
	l.readChar()
	return tok
}

func newToken(tokenType token.TypeToken, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch >= '9'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readChar() {
	defer func() {
		// next position is current position now
		l.position = l.readPosition
		// point to the next char
		l.readPosition++
	}()
	// Handling possible errors
	if l.readPosition >= len(l.input) {
		l.ch = 0
		return
	}
	// read next position.
	l.ch = l.input[l.readPosition]
}

// TestNextToken ...
func TestNextToken() {
	testInput := "=+(){},;"
	tests := []struct {
		expectedType    token.TypeToken
		expectedLiteral string
	}{

		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(testInput)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			fmt.Printf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
			return
		}

		if tok.Literal != tt.expectedLiteral {
			fmt.Printf("tests[%d] - literals wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
			return
		}
	}

	fmt.Print("Passed!\n")

}
