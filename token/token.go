package token

type Token struct {
	Type    TokenType
	Literal string
}

type TokenType string

// Specifying constant token types.
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers & literals.
	IDENT = "IDENT" // Identifier token.
	INT   = "INT"   // Integer type.

	// Operators.
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	BANG     = "!"
	LT       = "<"
	GT       = ">"
	EQUALS   = "=="
	NEQUALS  = "!="

	// Delimiters.
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords.
	FUNCTION = "FUNCTION"
	DECLARE  = "DECLARE"
	RETURN   = "RETURN"
	IF       = "IF"
	ELSE     = "ELSE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
)

var keywords = map[string]TokenType{
	"fn":    FUNCTION,
	"as":    DECLARE,
	"ret":   RETURN,
	"if":    IF,
	"else":  ELSE,
	"true":  TRUE,
	"false": FALSE,
}

// LookupIdent looks for an identifier and if it's a keyword, return it's representation.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
