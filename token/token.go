package token

// Type represents a type of token
type Type string

const (
	// ILLEGAL represents an illegal token
	ILLEGAL Type = "ILLEGAL"

	// EOF represents an end-of-file token
	EOF Type = "EOF"

	// IDENT represents an identifier e.g., add, foobar, x, y, ...
	IDENT Type = "IDENT"

	// INT represents an integer
	INT Type = "INT"

	// FLOAT represents a floating-point number
	FLOAT Type = "FLOAT"

	// ASSIGN is an assignment token
	ASSIGN Type = "="

	// PLUS is an addition token
	PLUS Type = "+"

	// MINUS is substration token
	MINUS Type = "-"

	// BANG is a bang token
	BANG Type = "!"

	// ASTERISK is a multiplication token
	ASTERISK Type = "*"

	// SLASH is a slash token
	SLASH Type = "/"

	// LT represents lesser than token
	LT Type = "<"

	// GT represents greater than token
	GT Type = ">"

	// EQ represents equals token
	EQ Type = "=="

	// NOTEQ represents not equals token
	NOTEQ Type = "!="

	// COMMA is a comma token
	COMMA Type = ","

	// SEMICOLON is a semicolon token
	SEMICOLON Type = ";"

	// LPAREN is a left parentheses token
	LPAREN Type = "("

	// RPAREN is a right parentheses token
	RPAREN Type = ")"

	// LBRACE is a left curly braces token
	LBRACE Type = "{"

	// RBRACE is a right curly braces token
	RBRACE Type = "}"

	// FUNCTION is a function token
	FUNCTION Type = "FUNCTION"

	// TRUE is a truthy boolean token
	TRUE Type = "TRUE"

	// FALSE is a falsy boolean token
	FALSE Type = "FALSE"

	// NULL represents a null value
	NULL Type = "NULL"

	// IF is an if statement token
	IF Type = "IF"

	// ELSE is an else statement token
	ELSE Type = "ELSE"

	// RETURN is a return statement token
	RETURN Type = "RETURN"

	// STRING represents a string literal
	STRING Type = "STRING"

	// LBRACKET is the left bracket token
	LBRACKET Type = "["

	// RBRACKET is the right bracket token
	RBRACKET Type = "]"

	// COLON is a colon token
	COLON Type = ":"

	// HASH is a hash token
	HASH Type = "#"

	// DOT is a dot token
	DOT Type = "."
)

// Token represents a single token
type Token struct {
	Type    Type
	Literal string
}

// keywords map are the supported language keywords
var keywords = map[string]Type{
	"fn":     FUNCTION,
	"true":   TRUE,
	"false":  FALSE,
	"null":   NULL,
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
