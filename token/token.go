package token

const (
	// ILLEGAL represents an illegal token
	ILLEGAL = "ILLEGAL"

	// EOF represents end-of-file token
	EOF = "EOF"

	// IDENT represents an identifier e.g add, foobar, x, y, ...
	IDENT = "IDENT"

	// INT represents an integer
	INT = "INT"

	// ASSIGN is assignment token
	ASSIGN = "="

	// PLUS is addition token
	PLUS = "+"

	// MINUS is substration token
	MINUS = "-"

	// BANG is a bang token
	BANG = "!"

	// ASTERISK is multiplication token
	ASTERISK = "*"

	// SLASH is a slash token
	SLASH = "/"

	// LT represents lesser than token
	LT = "<"

	// GT represents greater than token
	GT = ">"

	// EQ represents equals token
	EQ = "=="

	// NOTEQ represents not equals token
	NOTEQ = "!="

	// COMMA is a comma token
	COMMA = ","

	// SEMICOLON is a semicolon token
	SEMICOLON = ";"

	// LPAREN is a left parentheses token
	LPAREN = "("

	// RPAREN is a right parentheses token
	RPAREN = ")"

	// LBRACE is a left curly braces token
	LBRACE = "{"

	// RBRACE is a right curly braces token
	RBRACE = "}"

	// FUNCTION is a function token
	FUNCTION = "FUNCTION"

	// LET is a let token
	LET = "LET"

	// TRUE is a truthy boolean token
	TRUE = "TRUE"

	// FALSE is a falsy boolean token
	FALSE = "FALSE"

	// IF is a if statement token
	IF = "IF"

	// ELSE is a else statement token
	ELSE = "ELSE"

	// RETURN is a return statement token
	RETURN = "RETURN"

	STRING = "STRING"

	LBRACKET = "["

	RBRACKET = "]"
)

// Type represents type of a token
type Type string

// Token represents a single token
type Token struct {
	Type    Type
	Literal string
}

// keywords map are the supported language keywords
var keywords = map[string]Type{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdent checks if a string is an identifier
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
