package ast

import (
	"strings"
	"xlang/token"
)

// Node for the AST
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement ...
type Statement interface {
	Node
	statementNode()
}

// Expression produces a value
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of every AST
type Program struct {
	Statements []Statement
}

// TokenLiteral Root
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// String returns the program as a strings
func (p *Program) String() string {
	var str strings.Builder = strings.Builder{}
	for _, s := range p.Statements {
		str.WriteString(s.String() + "\n")
	}
	return str.String()
}

// Identifier is a variable name
type Identifier struct {
	Token token.Token // IDENT
	Value string      // value
}

func (i *Identifier) expressionNode() {}

// TokenLiteral .
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

func (i *Identifier) String() string { return i.Value }

// LetStatement represents things like let x = 5;
type LetStatement struct {
	Token token.Token // let
	Name  *Identifier // name
	Value Expression  // exp
}

func (ls *LetStatement) statementNode() {}

// TokenLiteral .
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	out := strings.Builder{}

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ReturnStatement represents a return statement duh
type ReturnStatement struct {
	Token       token.Token // the 'return'
	ReturnValue Expression  // what is returning
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral .
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	out := strings.Builder{}
	out.WriteString("return ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement represents things like x + 10;
type ExpressionStatement struct {
	Token      token.Token // the first token of the exp.
	Expression Expression
}

func (exst *ExpressionStatement) statementNode() {}

// TokenLiteral .
func (exst *ExpressionStatement) TokenLiteral() string { return exst.Token.Literal }

func (exst *ExpressionStatement) String() string {
	if exst.Expression != nil {
		return exst.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral ..
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

func (il *IntegerLiteral) String() string { return il.Token.Literal }

// PrefixExpression like -165 or !true
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral ..
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

func (pe *PrefixExpression) String() string {
	out := strings.Builder{}

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}
