package lexer

import (
	"bytes"
	"monkey/token"
)

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
	case '#':
		tok.Type = token.COMMENT
		tok.Literal = l.readLine()
	case '.':
		tok = l.newToken(token.DOT)
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
	case '[':
		tok = l.newToken(token.LBRACKET)
	case ']':
		tok = l.newToken(token.RBRACKET)
	case '\'':
		fallthrough
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case ':':
		tok = l.newToken(token.COLON)
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
			tok.Literal = l.readNumber()
			if bytes.Contains([]byte(tok.Literal), []byte(".")) ||
				bytes.Contains([]byte(tok.Literal), []byte("e")) ||
				bytes.Contains([]byte(tok.Literal), []byte("E")) {
				tok.Type = token.FLOAT
			} else {
				tok.Type = token.INT
			}

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

func (l *Lexer) readLine() string {
	position := l.position + 1

	for {
		l.readChar()

		if l.ch == '\r' || l.ch == '\n' || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	ch := l.ch
	var b bytes.Buffer

	for {
		l.readChar()

		// Support some basic escapes like \"
		if l.ch == '\\' {
			switch l.peekChar() {
			case '"':
				b.WriteByte('"')
			case '\'':
				b.WriteByte('\'')
			case 'n':
				b.WriteByte('\n')
			case 'r':
				b.WriteByte('\r')
			case 't':
				b.WriteByte('\t')
			case '\\':
				b.WriteByte('\\')
			}

			// Skip over the '\\' and the matched single escape char
			l.readChar()

			continue
		}

		if l.ch == ch || l.ch == 0 {
			break
		}

		b.WriteByte(l.ch)
	}

	return b.String()
}

func (l *Lexer) readNumber() string {
	position := l.position
	end := position

	for end < len(l.input) && l.isDigitAt(l.input[end]) {
		end++
	}

	if end+1 < len(l.input) && l.input[end] == '.' && l.isDigitAt(l.input[end+1]) {
		end += 2
		for end < len(l.input) && l.isDigitAt(l.input[end]) {
			end++
		}
	}

	if end < len(l.input) && (l.input[end] == 'e' || l.input[end] == 'E') {
		exponentEnd := end + 1
		if exponentEnd < len(l.input) && (l.input[exponentEnd] == '+' || l.input[exponentEnd] == '-') {
			exponentEnd++
		}
		if exponentEnd < len(l.input) && l.isDigitAt(l.input[exponentEnd]) {
			end = exponentEnd + 1
			for end < len(l.input) && l.isDigitAt(l.input[end]) {
				end++
			}
		}
	}

	l.position = end
	l.readPosition = end + 1
	if end >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[end]
	}

	return l.input[position:end]
}

func (l *Lexer) isDigitAt(ch byte) bool {
	return '0' <= ch && ch <= '9'
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
