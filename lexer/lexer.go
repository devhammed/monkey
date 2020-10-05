package lexer

import "monkey/token"

// Lexer is a source code lexing struct
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

// New creates a new instance of Lexer
func New(input string) *Lexer {
	l := &Lexer{input: input}

	l.readChar()

	return l
}

func (l *Lexer) newToken(tokenType token.Type) token.Token {
	return token.Token{Type: tokenType, Literal: string(l.ch)}
}

// NextToken returns the next character in source code stream
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = l.newToken(token.ASSIGN)
		}
	case '+':
		tok = l.newToken(token.PLUS)
	case '-':
		tok = l.newToken(token.MINUS)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOTEQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = l.newToken(token.BANG)
		}
	case '/':
		tok = l.newToken(token.SLASH)
	case '*':
		tok = l.newToken(token.ASTERISK)
	case '<':
		tok = l.newToken(token.LT)
	case '>':
		tok = l.newToken(token.GT)
	case ';':
		tok = l.newToken(token.SEMICOLON)
	case ',':
		tok = l.newToken(token.COMMA)
	case '(':
		tok = l.newToken(token.LPAREN)
	case ')':
		tok = l.newToken(token.RPAREN)
	case '{':
		tok = l.newToken(token.LBRACE)
	case '}':
		tok = l.newToken(token.RBRACE)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if l.isLetter() {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		}

		if l.isDigit() {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		}

		tok = l.newToken(token.ILLEGAL)
	}

	l.readChar()

	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	position := l.position

	for l.isDigit() {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for l.isLetter() {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) isDigit() bool {
	return '0' <= l.ch && l.ch <= '9'
}

func (l *Lexer) isLetter() bool {
	return 'a' <= l.ch && l.ch <= 'z' || 'A' <= l.ch && l.ch <= 'Z' || l.ch == '_'
}
