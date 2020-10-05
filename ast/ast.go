package ast

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
