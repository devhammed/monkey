package ast

import (
	"monkey/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &AssignmentExpression{
		Token: token.Token{Type: token.ASSIGN, Literal: "="},
		Left: &Identifier{
			Token: token.Token{Type: token.IDENT, Literal: "myVar"},
			Value: "myVar",
		},
		Value: &Identifier{
			Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
			Value: "anotherVar",
		},
	}

	programString := program.String()

	if programString != "myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", programString)
	}
}
