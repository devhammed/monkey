package ast

import "monkey/token"

// Node interface represents a single node
type Node interface {
	TokenLiteral() string
}

// Statement interface represents a statement
type Statement interface {
	Node
	statementNode()
}

// Expression interface represents an expression
type Expression interface {
	Node
	expressionNode()
}

// Program struct represents a source code
type Program struct {
	Statements []Statement
}

// TokenLiteral returns the first statement literal token
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
