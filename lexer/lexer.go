package lexer

import (
	"github.com/adrianplavka/fe/token"
)

// Lexer parses through input to look for tokens.
type Lexer struct {
	input        string
	ch           byte // Current char in examination.
	position     int  // Current position in input 			(points to current char).
	peekPosition int  // Current peaking position in input 	(after current char).
}

// New returns a Lexer based on it's input.
func New(input string) *Lexer {
	lex := &Lexer{input: input}
	lex.readChar()
	return lex
}

func (lex *Lexer) readChar() {
	if lex.peekPosition >= len(lex.input) {
		lex.ch = 0
	} else {
		lex.ch = lex.input[lex.peekPosition]
	}
	lex.position = lex.peekPosition
	lex.peekPosition++
}

// NextToken computes the next token based on the current char.
func (lex *Lexer) NextToken() token.Token {
	var tok token.Token

	lex.skipWhitespace()
	switch lex.ch {
	// A case for ASSIGN or EQUALS token.
	case '=':
		if lex.peekChar() == '=' {
			ch := lex.ch
			lex.readChar()
			tok = token.Token{Type: token.EQUALS, Literal: string(ch) + string(lex.ch)}
		} else {
			tok = newToken(token.ASSIGN, lex.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, lex.ch)
	case '(':
		tok = newToken(token.LPAREN, lex.ch)
	case ')':
		tok = newToken(token.RPAREN, lex.ch)
	case '{':
		tok = newToken(token.LBRACE, lex.ch)
	case '}':
		tok = newToken(token.RBRACE, lex.ch)
	case ',':
		tok = newToken(token.COMMA, lex.ch)
	case '+':
		tok = newToken(token.PLUS, lex.ch)
	case '-':
		tok = newToken(token.MINUS, lex.ch)
	case '*':
		tok = newToken(token.ASTERISK, lex.ch)
	case '/':
		tok = newToken(token.SLASH, lex.ch)
	// A case of BANG or NEQUALS token.
	case '!':
		if lex.peekChar() == '=' {
			ch := lex.ch
			lex.readChar()
			tok = token.Token{Type: token.NEQUALS, Literal: string(ch) + string(lex.ch)}
		} else {
			tok = newToken(token.BANG, lex.ch)
		}
	case '<':
		tok = newToken(token.LT, lex.ch)
	case '>':
		tok = newToken(token.GT, lex.ch)
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	// Default switch for all keywords, identifiers and illegal tokens.
	default:
		if isLetter(lex.ch) {
			tok.Literal = lex.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(lex.ch) {
			tok.Literal = lex.readDigit()
			tok.Type = token.INT
			return tok
		}
		// If it's not a letter we know, we return an ILLEGAL token.
		tok = newToken(token.ILLEGAL, lex.ch)
	}
	lex.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// Skips the whitespace in the lexing part.
func (lex *Lexer) skipWhitespace() {
	for lex.ch == ' ' || lex.ch == '\t' || lex.ch == '\n' || lex.ch == '\r' {
		lex.readChar()
	}
}

// Reads the identifier if it's a letter.
func (lex *Lexer) readIdentifier() string {
	position := lex.position
	for isLetter(lex.ch) {
		lex.readChar()
	}
	return lex.input[position:lex.position]
}

// Reads the identifier if it's a digit.
func (lex *Lexer) readDigit() string {
	position := lex.position
	for isDigit(lex.ch) {
		lex.readChar()
	}
	return lex.input[position:lex.position]
}

// Helper function for declaring a range of letters.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// Helper function for declaring a range of digits.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// Peeks the next char at the current char.
// If there is no such char, we return 0 (EOF).
func (lex *Lexer) peekChar() byte {
	if lex.peekPosition >= len(lex.input) {
		return 0
	}
	return lex.input[lex.peekPosition]
}
